import os
import subprocess

chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"
fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"

svg_path = os.path.join(fig_dir, "diag_01_nat_traversal.svg")
html_path = os.path.join(fig_dir, "diag_01_nat_traversal_hi.html")
png_path = os.path.join(fig_dir, "diag_01_nat_traversal.png")

with open(svg_path, "r", encoding="utf-8") as f:
    svg_content = f.read()

# Render at 3x resolution: 3600 x 1560 (Print quality ~400 DPI)
w, h = 3600, 1560

html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  * {{ margin: 0; padding: 0; box-sizing: border-box; }}
  html, body {{ width: {w}px; height: {h}px; background: #ffffff; overflow: hidden; }}
  svg {{ width: {w}px; height: {h}px; display: block; }}
</style>
</head>
<body>
  {svg_content}
</body>
</html>"""

with open(html_path, "w", encoding="utf-8") as f:
    f.write(html_content)

args = [
    chrome_path,
    "--headless=old",
    "--disable-gpu",
    "--hide-scrollbars",
    f"--window-size={w},{h}",
    f"--screenshot={png_path}",
    html_path
]

subprocess.run(args, capture_output=True)
print(f"Rendered High-Res PNG ({w}x{h}): {os.path.getsize(png_path)} bytes")
