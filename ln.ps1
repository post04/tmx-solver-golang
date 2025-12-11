# CountLines.ps1
$goFiles = Get-ChildItem -Path . -Recurse -Include *.go
$totalLines = 0

foreach ($file in $goFiles) {
    $lineCount = (Get-Content $file.FullName | Measure-Object -Line).Lines
    $totalLines += $lineCount
    Write-Output "$($file.FullName): $lineCount lines"
}

Write-Output "Total lines of Go code: $totalLines"