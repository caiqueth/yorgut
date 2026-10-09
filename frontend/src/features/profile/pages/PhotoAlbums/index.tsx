import { Box, Dialog, DialogContent, Stack, TextField, Typography } from "@mui/material";
import { useDeletePhoto, usePhotos, useUploadPhoto } from "api/photo/queries";
import { PhotoItem } from "api/generated/model";
import { Button } from "components/Button";
import { BoxWithTitle } from "components/BoxWithTitle";
import { dayjsFromDate } from "helpers/formatting/date";
import { uploadUrl } from "helpers/uploadUrl";
import { PICTURE_RULES_TEXT, validatePicture } from "helpers/validatePicture";
import React, { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";

// Same limit as the backend (service/photo.go).
const MAX_CAPTION_LENGTH = 200;

interface IProps { }
export const UserPhotoAlbumsList: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  const { data, isLoading } = usePhotos(userId);
  const { mutate: uploadPhoto, isPending: isUploading, isError: uploadFailed, reset: resetUpload } = useUploadPhoto(userId);
  const { mutate: deletePhoto, isPending: isDeleting, variables: deletingId } = useDeletePhoto(userId);

  const [picture, setPicture] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [caption, setCaption] = useState("");
  const [validationError, setValidationError] = useState<string | null>(null);
  const [opened, setOpened] = useState<PhotoItem | null>(null);
  const fileInput = useRef<HTMLInputElement>(null);

  // Release the local preview URL when it's replaced or the page closes.
  useEffect(() => () => { if (preview) URL.revokeObjectURL(preview); }, [preview]);

  const count = data?.count ?? 0;
  const maxPhotos = data?.maxPhotos ?? 15;
  const isFull = count >= maxPhotos;

  const clearPicked = () => {
    setPicture(null);
    setPreview(null);
    setCaption("");
    if (fileInput.current) fileInput.current.value = "";
  };

  const handlePick = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] ?? null;
    resetUpload();
    setValidationError(null);
    setPicture(null);
    setPreview(null);
    if (!file) return;

    const error = await validatePicture(file);
    if (error) {
      setValidationError(error);
      if (fileInput.current) fileInput.current.value = "";
      return;
    }
    setPicture(file);
    setPreview(URL.createObjectURL(file));
  };

  const handleUpload = () => {
    if (!picture) return;
    uploadPhoto({ picture, caption }, { onSuccess: clearPicked });
  };

  const photos = data?.photos ?? [];

  return (
    <Stack direction={"column"} gap={1}>
      {data?.canUpload && (
        <BoxWithTitle title={"adicionar foto"}>
          <Stack direction={"column"} gap={1} sx={{ p: 1 }}>
            {isFull ? (
              <Typography variant={"body2"}>
                Seu álbum já tem {maxPhotos} fotos. Exclua uma para adicionar outra.
              </Typography>
            ) : (
              <>
                <Typography variant={"body2"}>{PICTURE_RULES_TEXT}</Typography>
                <input ref={fileInput} type={"file"} accept={"image/png,image/jpeg"} onChange={handlePick} />
                {preview && <img src={preview} height={120} style={{ alignSelf: "start" }} alt={"pré-visualização"} />}
                <TextField
                  size={"small"}
                  label={"legenda (opcional)"}
                  value={caption}
                  onChange={(e) => setCaption(e.target.value)}
                  slotProps={{ htmlInput: { maxLength: MAX_CAPTION_LENGTH } }}
                />
                <Box>
                  <Button onClick={handleUpload} disabled={!picture} loading={isUploading}>
                    enviar foto
                  </Button>
                </Box>
              </>
            )}
            {validationError && (
              <Typography variant={"body2"} color={"error"}>{validationError}</Typography>
            )}
            {uploadFailed && (
              <Typography variant={"body2"} color={"error"}>
                Não foi possível enviar a foto. Tente novamente.
              </Typography>
            )}
          </Stack>
        </BoxWithTitle>
      )}

      <BoxWithTitle title={`álbum (${count}/${maxPhotos})`} isLoading={isLoading}>
        {!isLoading && photos.length === 0 ? (
          <Box sx={{ p: 1 }}>
            <Typography>Nenhuma foto ainda.</Typography>
          </Box>
        ) : (
          <Box
            sx={{
              p: 1,
              display: "grid",
              gridTemplateColumns: "repeat(auto-fill, minmax(140px, 1fr))",
              gap: 1,
            }}
          >
            {photos.map((photo) => (
              <Stack key={photo.id} direction={"column"} gap={0.5}>
                <Box
                  component={"img"}
                  src={uploadUrl(photo.url)}
                  alt={photo.caption || "foto"}
                  onClick={() => setOpened(photo)}
                  sx={{
                    width: "100%",
                    aspectRatio: "1",
                    objectFit: "cover",
                    cursor: "pointer",
                    border: "1px solid #ccc",
                  }}
                />
                {photo.caption && (
                  <Typography variant={"body2"} sx={{ overflowWrap: "anywhere" }}>{photo.caption}</Typography>
                )}
                {data?.canUpload && (
                  <Box>
                    <Button
                      onClick={() => deletePhoto(photo.id)}
                      loading={isDeleting && deletingId === photo.id}
                    >
                      excluir
                    </Button>
                  </Box>
                )}
              </Stack>
            ))}
          </Box>
        )}
      </BoxWithTitle>

      <Dialog open={!!opened} onClose={() => setOpened(null)} maxWidth={"lg"}>
        {opened && (
          <DialogContent>
            <Stack direction={"column"} gap={1} alignItems={"center"}>
              <img
                src={uploadUrl(opened.url)}
                alt={opened.caption || "foto"}
                style={{ maxWidth: "100%", maxHeight: "80vh" }}
              />
              {opened.caption && <Typography>{opened.caption}</Typography>}
              <Typography variant={"body2"}>
                {dayjsFromDate(opened.createdAt, "D [de] MMMM [de] YYYY")}
              </Typography>
            </Stack>
          </DialogContent>
        )}
      </Dialog>
    </Stack>
  );
};
