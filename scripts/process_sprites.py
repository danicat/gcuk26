import os, glob
from PIL import Image

def is_background_color(r, g, b, a):
    # Magenta / Pink (#FF00FF or similar)
    if r > 180 and b > 180 and g < 120:
        return True
    # White / Light Gray / Checkerboard
    if abs(r - g) < 20 and abs(g - b) < 20 and abs(r - b) < 20 and (r > 120 or r < 30):
        return True
    return False

def process_file(path):
    if not os.path.exists(path):
        return
    img = Image.open(path).convert("RGBA")
    pix = img.load()
    w, h = img.size

    # Floodfill background from corners if needed
    visited = set()
    queue = [(0, 0), (w-1, 0), (0, h-1), (w-1, h-1)]

    for q in queue:
        if q[0] >= 0 and q[0] < w and q[1] >= 0 and q[1] < h:
            r, g, b, a = pix[q[0], q[1]]
            if is_background_color(r, g, b, a):
                pix[q[0], q[1]] = (0, 0, 0, 0)

    # General chroma key scan
    for y in range(h):
        for x in range(w):
            r, g, b, a = pix[x, y]
            if is_background_color(r, g, b, a):
                pix[x, y] = (0, 0, 0, 0)

    # Tight crop around subject
    bbox = img.getbbox()
    if bbox:
        img = img.crop(bbox)

    img.save(path, "PNG")
    print(f"Processed {path}: final size {img.size}")

if __name__ == "__main__":
    for f in glob.glob("assets/*.png"):
        process_file(f)
