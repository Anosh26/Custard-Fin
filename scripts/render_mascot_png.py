import math
from PIL import Image, ImageDraw

def render_custard_mascot(out_path="docs/presentation/custard_mascot.png", size=512):
    # Create RGBA canvas with supersampling for crisp edges
    scale = 2
    canvas_size = size * scale
    im = Image.new("RGBA", (canvas_size, canvas_size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(im)

    def s(v):
        return v * (canvas_size / 200.0)

    # 1. Shadow / Glow under plate
    for r in range(12, 0, -2):
        alpha = int(25 * (1 - r / 12))
        draw.ellipse([s(100 - 70 - r), s(172 - 12 - r/2), s(100 + 70 + r), s(172 + 12 + r/2)], fill=(245, 158, 11, alpha))

    # 2. Plate rim & body
    draw.ellipse([s(20), s(152), s(180), s(184)], fill=(148, 163, 184, 255))
    draw.ellipse([s(22), s(151), s(178), s(181)], fill=(226, 232, 240, 255))
    draw.ellipse([s(32), s(155), s(168), s(177)], fill=(248, 250, 252, 255))
    draw.ellipse([s(44), s(157), s(156), s(173)], fill=(241, 245, 249, 255))

    # 3. Pudding Flan Body (Custard Pastel Yellow)
    pudding_poly = [
        (s(46), s(160)),
        (s(54), s(110)),
        (s(64), s(68)),
        (s(136), s(68)),
        (s(146), s(110)),
        (s(154), s(160)),
    ]
    draw.polygon(pudding_poly, fill=(254, 240, 138, 255)) # Pastel custard top
    # Bottom curved base
    draw.ellipse([s(46), s(150), s(154), s(170)], fill=(250, 204, 21, 255))
    draw.polygon([(s(46), s(160)), (s(64), s(68)), (s(136), s(68)), (s(154), s(160))], fill=(253, 224, 71, 255))

    # 4. Warm Caramel Drips
    # Caramel base top ellipse
    draw.ellipse([s(60), s(60), s(140), s(76)], fill=(180, 83, 9, 255))
    # Drips
    drips = [
        # x, y, r
        (68, 96, 12),
        (82, 108, 14),
        (98, 96, 11),
        (114, 112, 15),
        (128, 96, 13),
        (140, 92, 10),
    ]
    for dx, dy, dr in drips:
        draw.ellipse([s(dx - dr/2), s(dy - dr/2), s(dx + dr/2), s(dy + dr/2)], fill=(217, 119, 6, 255))
        draw.ellipse([s(dx - dr/3), s(dy - dr/3), s(dx + dr/3), s(dy + dr/3)], fill=(245, 158, 11, 255))

    # Caramel cap fill
    draw.ellipse([s(62), s(62), s(138), s(74)], fill=(245, 158, 11, 255))
    draw.ellipse([s(70), s(64), s(130), s(72)], fill=(251, 191, 36, 255))
    # Gloss highlight
    draw.ellipse([s(76), s(65), s(105), s(69)], fill=(255, 255, 255, 180))

    # 5. Kawaii Expression
    # Left eye curve
    draw.arc([s(76), s(118), s(92), s(130)], start=190, end=350, fill=(69, 26, 3, 255), width=int(s(3.5)))
    # Right eye curve
    draw.arc([s(108), s(118), s(124), s(130)], start=190, end=350, fill=(69, 26, 3, 255), width=int(s(3.5)))
    # Mouth smile
    draw.arc([s(94), s(130), s(106), s(140)], start=10, end=170, fill=(69, 26, 3, 255), width=int(s(3)))
    # Rosy cheeks (Pastel pink blush)
    draw.ellipse([s(66), s(128), s(78), s(136)], fill=(251, 113, 133, 160))
    draw.ellipse([s(122), s(128), s(134), s(136)], fill=(251, 113, 133, 160))

    # 6. Glossy Red Cherry on Top
    # Cherry Body
    draw.ellipse([s(84), s(34), s(116), s(66)], fill=(239, 68, 68, 255))
    draw.ellipse([s(88), s(38), s(112), s(62)], fill=(220, 38, 38, 255))
    # Cherry highlight
    draw.ellipse([s(90), s(40), s(98), s(46)], fill=(255, 255, 255, 210))
    # Stem
    draw.arc([s(95), s(18), s(135), s(52)], start=170, end=290, fill=(21, 128, 61, 255), width=int(s(3.5)))
    # Mint leaf
    draw.polygon([(s(112), s(34)), (s(126), s(28)), (s(124), s(40)), (s(116), s(40))], fill=(34, 197, 94, 255))

    # 7. Twinkling Sparkles
    def draw_star(cx, cy, r, color):
        pts = []
        for i in range(8):
            rad = r if i % 2 == 0 else r * 0.4
            ang = i * math.pi / 4
            pts.append((s(cx + rad * math.cos(ang)), s(cy + rad * math.sin(ang))))
        draw.polygon(pts, fill=color)

    draw_star(32, 70, 10, (253, 224, 71, 230))
    draw_star(168, 92, 9, (186, 230, 253, 230))
    draw_star(156, 46, 7, (253, 224, 71, 210))

    # Resample with Lanczos for smooth anti-aliased result
    final_img = im.resize((size, size), Image.Resampling.LANCZOS)
    final_img.save(out_path, "PNG")
    print(f"[OK] Saved high-resolution custard mascot PNG: {out_path}")

if __name__ == "__main__":
    render_custard_mascot()
