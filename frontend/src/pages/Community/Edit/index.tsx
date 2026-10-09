import { Box, Stack, TextField, Typography } from "@mui/material";
import {
  useCommunity,
  useCommunityProfile,
  useUpdateCommunity,
  useUpdateCommunityPicture,
} from "api/community/queries";
import { UpdateCommunityRequest } from "api/generated/model";
import { useMe } from "api/me/queries";
import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { uploadUrl } from "helpers/uploadUrl";
import { PICTURE_RULES_TEXT, validatePicture } from "helpers/validatePicture";
import React, { useEffect, useRef, useState } from "react";
import { Navigate, useNavigate, useParams } from "react-router-dom";

interface IProps { }
export const EditCommunity: React.FC<IProps> = () => {
  const params = useParams();
  const communityId = parseInt(params.communityId ?? "");
  const navigate = useNavigate();

  const { data: me } = useMe();
  const { data: community, isLoading } = useCommunity(communityId);
  const { data: profile } = useCommunityProfile(communityId);
  const { mutate: updateCommunity, isPending: isSaving, isError: saveFailed } = useUpdateCommunity(communityId);
  const {
    mutate: updatePicture,
    isPending: isUploading,
    isError: uploadFailed,
    reset: resetUpload,
  } = useUpdateCommunityPicture(communityId);

  const [form, setForm] = useState<UpdateCommunityRequest>({ name: "", category: "", description: "" });
  const [picture, setPicture] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [pictureError, setPictureError] = useState<string | null>(null);
  const fileInput = useRef<HTMLInputElement>(null);

  // Fill the form once the current data arrives.
  useEffect(() => {
    if (!community || !profile) return;
    setForm({
      name: community.name ?? "",
      category: community.category ?? "",
      description: profile.description ?? "",
    });
  }, [community, profile]);

  // Release the local preview URL when it's replaced or the page closes.
  useEffect(() => () => { if (preview) URL.revokeObjectURL(preview); }, [preview]);

  // Only the owner can edit; the backend enforces it too.
  if (me && community && me.id !== community.ownerId) {
    return <Navigate to={`/community/${communityId}`} replace />;
  }

  const setField = (key: keyof UpdateCommunityRequest) =>
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
      setForm((current) => ({ ...current, [key]: e.target.value }));

  const handlePickPicture = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] ?? null;
    resetUpload();
    setPictureError(null);
    setPicture(null);
    setPreview(null);
    if (!file) return;

    const error = await validatePicture(file);
    if (error) {
      setPictureError(error);
      if (fileInput.current) fileInput.current.value = "";
      return;
    }
    setPicture(file);
    setPreview(URL.createObjectURL(file));
  };

  const handleUploadPicture = () => {
    if (!picture) return;
    updatePicture(picture, {
      onSuccess: () => {
        setPicture(null);
        setPreview(null);
        if (fileInput.current) fileInput.current.value = "";
      },
    });
  };

  const handleSave = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    updateCommunity(form, { onSuccess: () => navigate(`/community/${communityId}`) });
  };

  const currentPicture = preview ?? uploadUrl(profile?.profilePicture);

  return (
    <Stack direction={"column"} gap={2}>
      <BoxWithTitle title={"imagem da comunidade"}>
        <Stack direction={"row"} gap={2} alignItems={"center"} sx={{ p: 1 }}>
          {currentPicture && <img src={currentPicture} height={96} alt={"imagem da comunidade"} />}
          <Stack direction={"column"} gap={1}>
            <input ref={fileInput} type={"file"} accept={"image/png,image/jpeg"} onChange={handlePickPicture} />
            <Typography variant={"body2"}>{PICTURE_RULES_TEXT}</Typography>
            <Box>
              <Button onClick={handleUploadPicture} disabled={!picture} loading={isUploading}>
                enviar imagem
              </Button>
            </Box>
            {pictureError && (
              <Typography variant={"body2"} color={"error"}>{pictureError}</Typography>
            )}
            {uploadFailed && (
              <Typography variant={"body2"} color={"error"}>
                Não foi possível enviar a imagem. Tente novamente.
              </Typography>
            )}
          </Stack>
        </Stack>
      </BoxWithTitle>

      <form onSubmit={handleSave}>
        <BoxWithTitle title={"editar comunidade"} isLoading={isLoading}>
          <Stack direction={"column"} gap={2} sx={{ p: 2 }}>
            <TextField
              label={"nome"}
              size={"small"}
              required
              value={form.name}
              onChange={setField("name")}
              slotProps={{ htmlInput: { maxLength: 100 } }}
            />
            <TextField
              label={"categoria"}
              size={"small"}
              required
              value={form.category}
              onChange={setField("category")}
              slotProps={{ htmlInput: { maxLength: 100 } }}
            />
            <TextField
              label={"descrição"}
              size={"small"}
              multiline
              minRows={4}
              value={form.description}
              onChange={setField("description")}
              slotProps={{ htmlInput: { maxLength: 2048 } }}
            />
            {saveFailed && (
              <Typography variant={"body2"} color={"error"}>
                Não foi possível salvar. Confira os campos e tente novamente.
              </Typography>
            )}
            <Stack direction={"row"} gap={1}>
              <Button type={"submit"} loading={isSaving}>salvar</Button>
              <Button onClick={() => navigate(`/community/${communityId}`)}>cancelar</Button>
            </Stack>
          </Stack>
        </BoxWithTitle>
      </form>
    </Stack>
  );
};
