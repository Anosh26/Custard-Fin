import os
import subprocess

fig_dir = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
chrome_path = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

diagrams = [
    ("diag_01_nat_traversal", 3600, 1560),
    ("diag_02_architecture_overview", 3600, 1740),
    ("diag_03_ingest_pipeline", 3600, 1440),
    ("diag_04_heartbeat_liveness", 3600, 1440),
    ("diag_05_dual_path_shield", 3600, 1680),
    ("diag_06_cost_comparison", 3600, 1560),
    ("diag_07_nvenc_benchmark", 3600, 1440)
]

for name, w, h in diagrams:
    svg_file = os.path.join(fig_dir, f"{name}.svg")
    png_file = os.path.join(fig_dir, f"{name}.png")
    html_file = os.path.join(fig_dir, f"{name}_render.html")

    if not os.path.exists(svg_file):
        print(f"Skipping {svg_file} (not found)")
        continue

    with open(svg_file, "r", encoding="utf-8") as f:
        svg_content = f.read()

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

    with open(html_file, "w", encoding="utf-8") as f:
        f.write(html_content)

    args = [
        chrome_path,
        "--headless=old",
        "--disable-gpu",
        "--hide-scrollbars",
        f"--window-size={w},{h}",
        f"--screenshot={png_file}",
        html_file
    ]

    subprocess.run(args, capture_output=True)
    size_kb = os.path.getsize(png_file) / 1024
    print(f"Rendered {name}.png at {w}x{h} Ultra-HD -> {size_kb:.1f} KB")

    # Clean up temporary html
    if os.path.exists(html_file):
        os.remove(html_file)

print("\nAll diagrams successfully re-rendered in Ultra-HD!")
