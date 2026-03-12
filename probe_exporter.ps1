$url = "http://10.0.72.42/apisec/probe/stress_test/probe_exporter/probe_exporter.exe"
$output = "C:\probe_exporter.exe"

Write-Host "正在从 $url 下载 probe_exporter.exe..."
Invoke-WebRequest -Uri $url -OutFile $output

Write-Host "下载完成，正在后台运行 $output..."
Start-Process -FilePath $output -WindowStyle Hidden
# Start-Process -FilePath $output
