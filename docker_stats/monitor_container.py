#!/usr/bin/env python3
# monitor_container.py
import docker
import time
import csv
import sys
from datetime import datetime
import signal

def signal_handler(sig, frame):
    print('\nDeteniendo monitoreo...')
    sys.exit(0)

def convert_bytes_to_mb(bytes_str):
    """Convierte bytes a MB"""
    try:
        if 'kB' in bytes_str:
            return float(bytes_str.replace('kB', '')) / 1000
        elif 'MB' in bytes_str:
            return float(bytes_str.replace('MB', ''))
        elif 'GB' in bytes_str:
            return float(bytes_str.replace('GB', '')) * 1000
        else:
            # Asumir bytes
            return float(bytes_str) / 1048576  # bytes to MB
    except:
        return 0

def monitor_container(container_name):
    client = docker.from_env()
    
    try:
        container = client.containers.get(container_name)
    except docker.errors.NotFound:
        print(f"Contenedor '{container_name}' no encontrado")
        return
    
    # Crear archivo CSV
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    filename = f"container_resources_{timestamp}.csv"
    
    with open(filename, 'w', newline='') as csvfile:
        fieldnames = ['timestamp', 'cpu_percent', 'memory_usage_mb', 
                     'memory_limit_mb', 'memory_percent', 'disk_read_mb', 
                     'disk_write_mb', 'network_rx_mb', 'network_tx_mb']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        
        print(f"Monitoreando {container_name}...")
        print(f"Log: {filename}")
        print("Presiona Ctrl+C para detener\n")
        
        # Variables para calcular diferencias
        prev_stats = None
        
        while True:
            try:
                # Obtener stats
                stats = container.stats(stream=False)
                timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
                
                # Calcular CPU
                cpu_percent = 0
                if prev_stats:
                    cpu_delta = stats['cpu_stats']['cpu_usage']['total_usage'] - prev_stats['cpu_stats']['cpu_usage']['total_usage']
                    system_delta = stats['cpu_stats']['system_cpu_usage'] - prev_stats['cpu_stats']['system_cpu_usage']
                    if system_delta > 0:
                        cpu_percent = (cpu_delta / system_delta) * len(stats['cpu_stats']['cpu_usage']['percpu_usage']) * 100
                
                # Memoria
                memory_usage_mb = stats['memory_stats']['usage'] / 1048576  # bytes to MB
                memory_limit_mb = stats['memory_stats']['limit'] / 1048576
                memory_percent = (memory_usage_mb / memory_limit_mb) * 100
                
                # Disk I/O
                disk_read_mb = 0
                disk_write_mb = 0
                if 'blkio_stats' in stats and 'io_service_bytes_recursive' in stats['blkio_stats']:
                    for item in stats['blkio_stats']['io_service_bytes_recursive']:
                        if item['op'] == 'Read':
                            disk_read_mb += item['value'] / 1048576
                        elif item['op'] == 'Write':
                            disk_write_mb += item['value'] / 1048576
                
                # Network I/O
                network_rx_mb = 0
                network_tx_mb = 0
                if 'networks' in stats:
                    for interface in stats['networks'].values():
                        network_rx_mb += interface['rx_bytes'] / 1048576
                        network_tx_mb += interface['tx_bytes'] / 1048576
                
                # Escribir a CSV
                writer.writerow({
                    'timestamp': timestamp,
                    'cpu_percent': round(cpu_percent, 2),
                    'memory_usage_mb': round(memory_usage_mb, 2),
                    'memory_limit_mb': round(memory_limit_mb, 2),
                    'memory_percent': round(memory_percent, 2),
                    'disk_read_mb': round(disk_read_mb, 2),
                    'disk_write_mb': round(disk_write_mb, 2),
                    'network_rx_mb': round(network_rx_mb, 2),
                    'network_tx_mb': round(network_tx_mb, 2)
                })
                
                # Mostrar en consola
                print(f"[{timestamp}] CPU: {cpu_percent:.1f}% | RAM: {memory_usage_mb:.1f}MB ({memory_percent:.1f}%) | Disk R/W: {disk_read_mb:.1f}/{disk_write_mb:.1f}MB")
                
                csvfile.flush()  # Asegurar escritura inmediata
                prev_stats = stats
                
            except Exception as e:
                print(f"Error: {e}")
            
            time.sleep(1)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Uso: python3 monitor_container.py <nombre_del_contenedor>")
        sys.exit(1)
    
    signal.signal(signal.SIGINT, signal_handler)
    monitor_container(sys.argv[1])
