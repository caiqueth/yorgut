// Upload limits for every picture (album photos, profile and community
// pictures). The backend enforces the same ones (service/picture.go) and
// shrinks what it stores to 1024px and 1 MB; checking here gives instant
// feedback before uploading.
export const MAX_PICTURE_BYTES = 5 * 1024 * 1024;
export const MAX_PICTURE_PIXELS = 50_000_000;
export const PICTURE_RULES_TEXT =
  "PNG ou JPEG, até 5 MB. Fotos grandes são reduzidas automaticamente.";

// Resolves to an error message, or null when the file is acceptable.
export const validatePicture = async (file: File): Promise<string | null> => {
  if (!["image/png", "image/jpeg"].includes(file.type)) {
    return "A foto precisa ser PNG ou JPEG.";
  }
  if (file.size > MAX_PICTURE_BYTES) {
    return "A foto precisa ter no máximo 5 MB.";
  }
  try {
    const bitmap = await createImageBitmap(file);
    const { width, height } = bitmap;
    bitmap.close();
    if (width * height > MAX_PICTURE_PIXELS) {
      return `A foto tem ${width}x${height} pixels; o máximo é 50 megapixels.`;
    }
  } catch {
    return "Não foi possível ler a imagem.";
  }
  return null;
};
