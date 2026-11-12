import os
import time
from playwright.sync_api import sync_playwright

FIG_DIR = r"C:\Users\Anosh\Documents\Notes\Projects\Cloud_Computing\sem_end\Custard Fin\report\figures"
os.makedirs(FIG_DIR, exist_ok=True)

CHROME_PATH = r"C:\Program Files\Google\Chrome\Application\chrome.exe"
HOST = "http://demo.custardfin.local"
LOCAL_EDGE = "http://127.0.0.1:8090"

def capture_all():
    with sync_playwright() as p:
        browser = p.chromium.launch(
            executable_path=CHROME_PATH,
            headless=True,
            args=[
                "--disable-extensions",
                "--disable-web-security",
                "--host-resolver-rules=MAP demo.custardfin.local 8.234.112.169"
            ]
        )

        context = browser.new_context(
            viewport={"width": 1440, "height": 900},
            device_scale_factor=2 # 2880x1800 Retina resolution
        )

        page = context.new_page()

        # Helper to set auth & local edge URL
        def set_state():
            page.evaluate("""() => {
                localStorage.setItem("custardfin_user", JSON.stringify({
                    id: "71675f0d-1efa-44f4-af0f-a8ef44995a23",
                    username: "anosh"
                }));
                localStorage.setItem("custardfin_jwt_token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.demo");
                localStorage.setItem("custardfin_edge_api_url", "http://127.0.0.1:8090");
            }""")

        # 1. Catalog Feed (Live from GCP VM)
        print("1. Capturing ui_catalog_feed.png...")
        page.goto(HOST, wait_until="networkidle")
        set_state()
        page.reload(wait_until="networkidle")
        page.wait_for_timeout(2000)
        catalog_path = os.path.join(FIG_DIR, "ui_catalog_feed.png")
        page.screenshot(path=catalog_path)
        print(f"   -> Saved: {catalog_path} ({os.path.getsize(catalog_path)} bytes)")

        # 2. Player with Stats for Nerds HUD
        print("2. Capturing ui_player_stats.png...")
        page.goto(f"{HOST}/watch/41e79799-877e-4feb-9c26-6014d07d2bba", wait_until="networkidle")
        set_state()
        page.wait_for_timeout(2500)
        # Attempt to trigger video play
        page.evaluate("""() => {
            const v = document.querySelector('video');
            if (v) {
                v.muted = true;
                v.play().catch(() => {});
            }
        }""")
        # Open Stats for Nerds
        stats_btn = page.locator(".stats-toggle-btn")
        if stats_btn.count() > 0:
            stats_btn.click()
            page.wait_for_timeout(1000)
        player_path = os.path.join(FIG_DIR, "ui_player_stats.png")
        page.screenshot(path=player_path)
        print(f"   -> Saved: {player_path} ({os.path.getsize(player_path)} bytes)")

        # 3. Edge Library Modal (SQLite on this PC)
        print("3. Capturing ui_edge_library.png...")
        page.goto(HOST, wait_until="networkidle")
        set_state()
        page.reload(wait_until="networkidle")
        page.wait_for_timeout(1500)
        lib_btn = page.locator("button:has-text('My Library')")
        if lib_btn.count() > 0:
            lib_btn.click()
            page.wait_for_timeout(2500)
        lib_path = os.path.join(FIG_DIR, "ui_edge_library.png")
        page.screenshot(path=lib_path)
        print(f"   -> Saved: {lib_path} ({os.path.getsize(lib_path)} bytes)")

        # 4. Dev Console Modal
        print("4. Capturing ui_dev_settings.png...")
        page.goto(HOST, wait_until="networkidle")
        set_state()
        page.wait_for_timeout(1000)
        dev_btn = page.locator("button:has-text('Dev Console')")
        if dev_btn.count() > 0:
            dev_btn.click()
            page.wait_for_timeout(1500)
        dev_path = os.path.join(FIG_DIR, "ui_dev_settings.png")
        page.screenshot(path=dev_path)
        print(f"   -> Saved: {dev_path} ({os.path.getsize(dev_path)} bytes)")

        # 5. Standalone Mini-Jellyfin (Edge Local Streamer)
        print("5. Capturing ui_mini_jellyfin.png...")
        try:
            page.goto(LOCAL_EDGE, wait_until="networkidle", timeout=5000)
            page.wait_for_timeout(1500)
            mini_path = os.path.join(FIG_DIR, "ui_mini_jellyfin.png")
            page.screenshot(path=mini_path)
            print(f"   -> Saved: {mini_path} ({os.path.getsize(mini_path)} bytes)")
        except Exception as e:
            print(f"   -> Could not capture mini-jellyfin: {e}")

        browser.close()
        print("\nAll live screenshots captured successfully!")

if __name__ == "__main__":
    capture_all()