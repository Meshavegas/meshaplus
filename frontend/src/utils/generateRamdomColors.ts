
/**
 * Generates a random color (hex string) that ensures white text (#FFFFFF) is always visible on it.
 * This is achieved by generating only "dark enough" colors (low luminance).
 */
function generateRandomColorForWhiteText(): string {
  // Helper to calculate luminance of a color
  function luminance(r: number, g: number, b: number): number {
    // sRGB luminance(Y) values
    return 0.2126 * r + 0.7152 * g + 0.0722 * b
  }

  let r, g, b, lum
  do {
    r = Math.floor(Math.random() * 256)
    g = Math.floor(Math.random() * 256)
    b = Math.floor(Math.random() * 256)
    lum = luminance(r, g, b)
    // Luminance threshold: lower means darker. 128 is a good cutoff for white text.
  } while (lum > 128)

  // Convert to hex string
  const hex = (n: number) => n.toString(16).padStart(2, '0')
  return `#${hex(r)}${hex(g)}${hex(b)}`
}

export default generateRandomColorForWhiteText