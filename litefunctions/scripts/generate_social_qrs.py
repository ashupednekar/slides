#!/usr/bin/env python3
"""Generate social QR PNGs for thank-you slide (repo root thank-you.html). Requires: pip install qrcode[pil]"""

from pathlib import Path

import qrcode
from qrcode.constants import ERROR_CORRECT_M

# Same URLs as <a href> in ../../thank-you.html (must stay in sync)
URLS: dict[str, str] = {
    "linkedinqr.png": "https://www.linkedin.com/in/ashupednekar",
    "xqr.png": "https://x.com/ashupednekar49",
    "threadsqr.png": "https://www.threads.com/@ashupednekar",
    "blueskyqr.png": "https://bsky.app/profile/ashupednekar.bsky.social",
}


def main() -> None:
    out_dir = Path(__file__).resolve().parent.parent / "images"
    out_dir.mkdir(parents=True, exist_ok=True)

    for filename, data in URLS.items():
        qr = qrcode.QRCode(
            version=None,
            error_correction=ERROR_CORRECT_M,
            box_size=8,
            border=2,
        )
        qr.add_data(data)
        qr.make(fit=True)
        img = qr.make_image(fill_color="black", back_color="white")
        path = out_dir / filename
        img.save(path)
        print(f"Wrote {path} ({data!r})")


if __name__ == "__main__":
    main()
