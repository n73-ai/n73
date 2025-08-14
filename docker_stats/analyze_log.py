# analyze_log.py
import pandas as pd
import matplotlib.pyplot as plt

# Leer el CSV
df = pd.read_csv('c.csv')

# Estadísticas
print("=== RESUMEN ===")
print(f"CPU máximo: {df['cpu_percent'].max():.2f}%")
print(f"CPU promedio: {df['cpu_percent'].mean():.2f}%")
print(f"RAM máxima: {df['memory_usage_mb'].max():.2f}MB")
print(f"RAM promedio: {df['memory_usage_mb'].mean():.2f}MB")

# Gráfico
fig, axes = plt.subplots(2, 2, figsize=(15, 10))
df['cpu_percent'].plot(ax=axes[0,0], title='CPU Usage %')
df['memory_usage_mb'].plot(ax=axes[0,1], title='Memory Usage MB')
df['disk_read_mb'].plot(ax=axes[1,0], title='Disk Read MB')
df['network_rx_mb'].plot(ax=axes[1,1], title='Network RX MB')
plt.tight_layout()
plt.savefig('resource_usage.png')
