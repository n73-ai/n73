/*
// Enviar a usuario específico
SendToUser("user123", "Hello!")
// Broadcast a todos

BroadcastMessage("Server announcement")
// Ver estadísticas
fmt.Printf("Clients connected: %d", GetConnectedClientsCount())
fmt.Printf("User connections: %d", GetUserConnectionsCount("user123"))
*/
package handlers

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

// Client representa un cliente WebSocket conectado
type Client struct {
	isClosing bool
	mu        sync.Mutex
	ID        string
	conn      *websocket.Conn
	lastPing  time.Time
}

// Hub maneja todas las conexiones WebSocket
type Hub struct {
	clients    map[*websocket.Conn]*Client
	register   chan *Client
	unregister chan *websocket.Conn
	broadcast  chan []byte
	mu         sync.RWMutex
	running    bool
}

// hub es la instancia global del hub
var hub = &Hub{
	clients:    make(map[*websocket.Conn]*Client),
	register:   make(chan *Client, 256),
	unregister: make(chan *websocket.Conn, 256),
	broadcast:  make(chan []byte, 256),
}

// init inicializa el hub automáticamente
func init() {
	go hub.Run()
	go hub.startCleanupTicker()
}

// SendToUser envía un mensaje a un usuario específico por ID
func (h *Hub) SendToUser(id string, msg string) error {
	if id == "" || msg == "" {
		return nil
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for conn, client := range h.clients {
		if client.ID == id {
			client.mu.Lock()
			if client.isClosing {
				client.mu.Unlock()
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				client.isClosing = true
				client.mu.Unlock()
				// Cleanup asíncrono para evitar deadlock
				go func(c *websocket.Conn) {
					select {
					case h.unregister <- c:
					case <-time.After(1 * time.Second):
						log.Printf("Failed to unregister client after timeout")
					}
				}(conn)
				return err
			}
			client.mu.Unlock()
			return nil
		}
	}
	return nil
}

// BroadcastToAll envía un mensaje a todos los clientes conectados
func (h *Hub) BroadcastToAll(msg string) {
	if msg == "" {
		return
	}

	select {
	case h.broadcast <- []byte(msg):
	case <-time.After(1 * time.Second):
		log.Printf("Broadcast channel full, message dropped")
	}
}

// GetConnectedClients retorna el número de clientes conectados
func (h *Hub) GetConnectedClients() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetClientsByID retorna una lista de IDs de clientes conectados
func (h *Hub) GetClientsByID(id string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, client := range h.clients {
		if client.ID == id {
			count++
		}
	}
	return count
}

// Run ejecuta el loop principal del hub
func (h *Hub) Run() {
	h.running = true
	log.Println("WebSocket Hub started")

	for h.running {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case conn := <-h.unregister:
			h.unregisterClient(conn)
		}
	}
}

// registerClient registra un nuevo cliente
func (h *Hub) registerClient(client *Client) {
	if client == nil || client.conn == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.conn] = client
	client.lastPing = time.Now()

	log.Printf("Client registered: ID=%s, Total clients=%d", client.ID, len(h.clients))
}

// broadcastMessage envía un mensaje a todos los clientes
func (h *Hub) broadcastMessage(message []byte) {
	h.mu.RLock()
	clientsCopy := make(map[*websocket.Conn]*Client)
	for conn, client := range h.clients {
		clientsCopy[conn] = client
	}
	h.mu.RUnlock()

	for conn, client := range clientsCopy {
		client.mu.Lock()
		if client.isClosing {
			client.mu.Unlock()
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			client.isClosing = true
			client.mu.Unlock()
			// Cleanup asíncrono
			go func(c *websocket.Conn) {
				select {
				case h.unregister <- c:
				case <-time.After(1 * time.Second):
					log.Printf("Failed to unregister client after broadcast error")
				}
			}(conn)
			continue
		}
		client.mu.Unlock()
	}
}

// unregisterClient desregistra y limpia un cliente
func (h *Hub) unregisterClient(conn *websocket.Conn) {
	if conn == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if client, exists := h.clients[conn]; exists {
		client.mu.Lock()
		client.isClosing = true
		client.mu.Unlock()

		// Cleanup completo
		delete(h.clients, conn)
		conn.Close()

		log.Printf("Client unregistered: ID=%s, Total clients=%d", client.ID, len(h.clients))
	}
}

// startCleanupTicker inicia el limpiador periódico de conexiones muertas
func (h *Hub) startCleanupTicker() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.cleanupDeadConnections()
	}
}

// cleanupDeadConnections limpia conexiones muertas usando ping/pong
func (h *Hub) cleanupDeadConnections() {
	h.mu.RLock()
	clientsToCheck := make(map[*websocket.Conn]*Client)
	for conn, client := range h.clients {
		clientsToCheck[conn] = client
	}
	h.mu.RUnlock()

	for conn, client := range clientsToCheck {
		client.mu.Lock()
		if client.isClosing {
			client.mu.Unlock()
			continue
		}

		// Verificar si la conexión es muy antigua sin ping
		if time.Since(client.lastPing) > 2*time.Minute {
			client.isClosing = true
			client.mu.Unlock()
			go func(c *websocket.Conn) {
				select {
				case h.unregister <- c:
				case <-time.After(1 * time.Second):
					log.Printf("Failed to unregister stale client")
				}
			}(conn)
			continue
		}

		// Enviar ping para verificar conexión
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			client.isClosing = true
			client.mu.Unlock()
			go func(c *websocket.Conn) {
				select {
				case h.unregister <- c:
				case <-time.After(1 * time.Second):
					log.Printf("Failed to unregister client after ping error")
				}
			}(conn)
			continue
		}
		client.mu.Unlock()
	}
}

// Stop para el hub gracefully
func (h *Hub) Stop() {
	h.running = false
	log.Println("WebSocket Hub stopping...")

	// Cerrar todas las conexiones
	h.mu.Lock()
	for conn, client := range h.clients {
		client.mu.Lock()
		client.isClosing = true
		client.mu.Unlock()
		conn.Close()
	}
	h.clients = make(map[*websocket.Conn]*Client)
	h.mu.Unlock()

	log.Println("WebSocket Hub stopped")
}

// FeedChat maneja una nueva conexión WebSocket
func FeedChat(c *websocket.Conn) {
	// Validar parámetros
	id := c.Query("id")
	if id == "" {
		log.Printf("WebSocket connection rejected: missing id parameter")
		c.WriteMessage(websocket.CloseMessage, []byte("Missing id parameter"))
		c.Close()
		return
	}

	// Configurar timeouts
	c.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Crear cliente
	client := &Client{
		ID:       id,
		conn:     c,
		lastPing: time.Now(),
	}

	// Configurar handlers
	c.SetPongHandler(func(string) error {
		client.mu.Lock()
		client.lastPing = time.Now()
		client.mu.Unlock()
		return c.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	c.SetCloseHandler(func(code int, text string) error {
		log.Printf("WebSocket connection closed: ID=%s, Code=%d, Text=%s", id, code, text)
		return nil
	})

	// Registrar cliente
	select {
	case hub.register <- client:
	case <-time.After(5 * time.Second):
		log.Printf("Failed to register client: timeout")
		c.Close()
		return
	}

	// Cleanup garantizado
	defer func() {
		select {
		case hub.unregister <- c:
		case <-time.After(5 * time.Second):
			log.Printf("Failed to unregister client: timeout")
		}
	}()

	log.Printf("WebSocket connection established: ID=%s", id)

	// Loop para mantener conexión viva y leer mensajes
	for {
		// Leer mensajes (principalmente para detectar disconnects)
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Actualizar último ping en cualquier mensaje recibido
		client.mu.Lock()
		client.lastPing = time.Now()
		client.mu.Unlock()

		// Resetear read deadline
		c.SetReadDeadline(time.Now().Add(60 * time.Second))

		// Procesar mensajes si es necesario
		if messageType == websocket.TextMessage {
			log.Printf("Received message from client %s: %s", id, string(message))
			// Aquí puedes agregar lógica para procesar mensajes del cliente
		}
	}

	log.Printf("WebSocket connection ended: ID=%s", id)
}

// Funciones públicas para usar desde otros paquetes

// SendToUser envía un mensaje a un usuario específico
func SendToUser(id string, msg string) error {
	return hub.SendToUser(id, msg)
}

// BroadcastMessage envía un mensaje a todos los clientes conectados
func BroadcastMessage(msg string) {
	hub.BroadcastToAll(msg)
}

// GetConnectedClientsCount retorna el número de clientes conectados
func GetConnectedClientsCount() int {
	return hub.GetConnectedClients()
}

// GetUserConnectionsCount retorna cuántas conexiones tiene un usuario específico
func GetUserConnectionsCount(id string) int {
	return hub.GetClientsByID(id)
}

// StopHub para el hub (útil para tests o shutdown graceful)
func StopHub() {
	hub.Stop()
}

// RunHub inicia el hub manualmente (opcional, ya se inicia automáticamente)
func RunHub() {
	if !hub.running {
		go hub.Run()
		go hub.startCleanupTicker()
	}
}
