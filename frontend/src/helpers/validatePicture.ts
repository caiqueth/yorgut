// Upload limits for every picture (album photos and profile pictures). The
// backend enforces the same ones (service/picture.go); checking here gives
// instant feedback before uploading.
export const MAX_PICTURE_BYTES = 1024 * 1024;
export const MAX_PICTURE_DIMENSION = 1024;
export const PICTURE_RULES_TEXT =
  `PNG ou JPEG, até 1 MB e ${MAX_PICTURE_DIMENSION}x${MAX_PICTURE_DIMENSION} pixels.`;

// Resolves to an error message, or null when the file is acceptable.
export const validatePicture = async (file: File): Promise<string | null> => {
  if (!["image/png", "image/jpeg"].includes(file.type)) {
    return "A foto precisa ser PNG ou JPEG.";
  }
  if (file.size > MAX_PICTURE_BYTES) {
    return "A foto precisa ter no máximo 1 MB.";
  }
  try {
    const bitmap = await createImageBitmap(file);
    const { width, height } = bitmap;
    bitmap.close();
    if (width > MAX_PICTURE_DIMENSION || height > MAX_PICTURE_DIMENSION) {
      return `A foto tem ${width}x${height} pixels; o máximo é ${MAX_PICTURE_DIMENSION}x${MAX_PICTURE_DIMENSION}.`;
    }
  } catch {
    return "Não foi possível ler a imagem.";
  }
  return null;
};
