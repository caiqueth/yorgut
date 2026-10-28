import React from "react";
import { Box, Stack, Typography } from "@mui/material";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Button } from "components/Button";
import { PICTURE_RULES_TEXT } from "helpers/validatePicture";
import { usePicturePicker } from "hooks/usePicturePicker";

interface IPictureUploadBoxProps {
  // The box title, e.g. "foto do perfil".
  title: string;
  // What the picture is called in the button and messages: "enviar foto".
  noun: "foto" | "imagem";
  // The saved picture, shown until a new one is picked.
  currentPicture?: string | null;
  isUploading: boolean;
  // The last upload's error, with the API's message; null when it didn't fail.
  uploadError: Error | null;
  // Sends the picture; onSent clears the picker once it's saved.
  onUpload: (picture: File, onSent: () => void) => void;
  // Called when a file is picked, e.g. to clear the last upload error.
  onPick?: () => void;
}

// Replaces a profile's or community's picture: pick, preview, send.
export const PictureUploadBox: React.FC<IPictureUploadBoxProps> = ({
  title,
  noun,
  currentPicture,
  isUploading,
  uploadError,
  onUpload,
  onPick,
}) => {
  const { inputRef, picture, preview, error, pick, clear } = usePicturePicker();
  const shown = preview ?? currentPicture ?? undefined;

  const handlePick = (e: React.ChangeEvent<HTMLInputElement>) => {
    onPick?.();
    pick(e);
  };

  return (
    <BoxWithTitle title={title}>
      {/* Picture next to the controls; under each other on phones. */}
      <Stack direction={{ xs: "column", sm: "row" }} gap={2} alignItems={{ xs: "flex-start", sm: "center" }} sx={{ p: 1 }}>
        {shown && <img src={shown} height={96} alt={title} />}
        <Stack direction={"column"} gap={1}>
          <input ref={inputRef} type={"file"} style={{ maxWidth: "100%" }} accept={"image/png,image/jpeg"} onChange={handlePick} />
          <Typography variant={"body2"}>{PICTURE_RULES_TEXT}</Typography>
          <Box>
            <Button onClick={() => picture && onUpload(picture, clear)} disabled={!picture} loading={isUploading}>
              enviar {noun}
            </Button>
          </Box>
          {error && (
            <Typography variant={"body2"} color={"error"}>{error}</Typography>
          )}
          {uploadError && (
            <Typography variant={"body2"} color={"error"}>
              {uploadError.message}
            </Typography>
          )}
        </Stack>
      </Stack>
    </BoxWithTitle>
  );
};
