# 🎨 Modern Pixel Paint

A simple, modern pixel paint application built with [Ebiten](https://ebiten.org/) in Go.  
Features a fixed top bar with tool icons, a color palette, and a brush size slider on the right.

---

## ✨ Features

- 🖌️ **Brush & Eraser:** Switch between brush and eraser modes.
- 🗑️ **Clear & Save:** Clear the canvas or save your artwork as `output.png`.
- 🌈 **Color Palette:** Select from a range of preset colors.
- 📏 **Brush Size Slider:** Adjust brush size using the slider.
- 🟪 **Live Cursor Preview:** See a square preview of your brush under the cursor.
- 🖼️ **Resizable Window:** The canvas adapts to window resizing.

---

## 🎮 Controls

| Key | Action         |
|-----|---------------|
| **B** | Select Brush   |
| **E** | Select Eraser  |
| **C** | Clear Canvas   |
| **S** | Save Canvas    |
| **Mouse** | Draw, select colors, adjust brush size |

---

## 🚀 How to Run

1. **Install Go:** [Download Go](https://golang.org/dl/)
2. **Install Ebiten:**
    ```sh
    go get github.com/hajimehoshi/ebiten/v2
    ```
3. **Run the app:**
    ```sh
    go run main.go
    ```

---

## 📁 File Structure

- `main.go` – Main application code

---

## 🖼️ Screenshot

![screenshot](screenshot.png) <!-- Add a screenshot if available -->

---

## 📜 License

MIT License

---

Made with ❤️ using
