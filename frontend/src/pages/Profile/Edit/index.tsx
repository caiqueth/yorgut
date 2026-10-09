import { Box, Stack, TextField, Typography } from "@mui/material";
import { useMe, useMeProfile, useUpdateMe, useUpdateProfilePicture } from "api/me/queries";
import { UpdateMeRequest } from "api/generated/model";
import { Button } from "components/atoms/Button";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { uploadUrl } from "helpers/uploadUrl";
import { PICTURE_RULES_TEXT, validatePicture } from "helpers/validatePicture";
import React, { useEffect, useState } from "react";
import { Navigate, useNavigate, useParams } from "react-router-dom";
import { profileFields } from "../profileFields";

const emptyForm: UpdateMeRequest = {
  firstName: "",
  lastName: "",
  birthDate: "",
  ...Object.fromEntries(profileFields.map(({ key }) => [key, ""])),
} as UpdateMeRequest;

// The API sends birth dates as "1990-03-15T00:00:00-03:00". Users created on
// first login get Go's zero time ("0001-01-01..."), shown as empty.
const toDateInput = (birthDate?: string) =>
  birthDate && !birthDate.startsWith("0001-") ? birthDate.slice(0, 10) : "";

interface IProps { }
export const EditProfile: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");
  const navigate = useNavigate();

  const { data: me } = useMe();
  const { data: profile, isLoading } = useMeProfile();
  const { mutate: updateMe, isPending: isSaving, error: saveError } = useUpdateMe(userId);
  const {
    mutate: updatePicture,
    isPending: isUploading,
    error: uploadError,
    reset: resetUpload,
  } = useUpdateProfilePicture(userId);

  const [form, setForm] = useState<UpdateMeRequest>(emptyForm);
  const [picture, setPicture] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [pictureError, setPictureError] = useState<string | null>(null);

  // Fill the form once the current data arrives.
  useEffect(() => {
    if (!profile) return;
    setForm({
      firstName: profile.firstName ?? "",
      lastName: profile.lastName ?? "",
      birthDate: toDateInput(profile.birthDate),
      ...Object.fromEntries(profileFields.map(({ key }) => [key, profile[key] ?? ""])),
    } as UpdateMeRequest);
  }, [profile]);

  // Release the local preview URL when it's replaced or the page closes.
  useEffect(() => () => { if (preview) URL.revokeObjectURL(preview); }, [preview]);

  if (me && me.id !== userId) {
    return <Navigate to={`/profile/${userId}`} replace />;
  }

  const setField = (key: keyof UpdateMeRequest) =>
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
      setForm((current) => ({ ...current, [key]: e.target.value }));

  const handlePickPicture = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const input = e.target;
    const file = input.files?.[0] ?? null;
    resetUpload();
    setPictureError(null);
    setPicture(null);
    setPreview(null);
    if (!file) return;

    const error = await validatePicture(file);
    if (error) {
      setPictureError(error);
      input.value = "";
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
      },
    });
  };

  const handleSave = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    updateMe(form, { onSuccess: () => navigate(`/profile/${userId}`) });
  };

  const currentPicture = preview ?? uploadUrl(profile?.profilePicture);

  return (
    <Stack direction={"column"} gap={2}>
      <BoxWithTitle title={"foto do perfil"}>
        <Stack direction={"row"} gap={2} alignItems={"center"} sx={{ p: 1 }}>
          {currentPicture && <img src={currentPicture} height={96} alt={"foto do perfil"} />}
          <Stack direction={"column"} gap={1}>
            <input type={"file"} accept={"image/png,image/jpeg"} onChange={handlePickPicture} />
            <Typography variant={"body2"}>{PICTURE_RULES_TEXT}</Typography>
            <Box>
              <Button onClick={handleUploadPicture} disabled={!picture} loading={isUploading}>
                enviar foto
              </Button>
            </Box>
            {pictureError && (
              <Typography variant={"body2"} color={"error"}>{pictureError}</Typography>
            )}
            {uploadError && (
              <Typography variant={"body2"} color={"error"}>
                Não foi possível enviar a foto. Tente novamente.
              </Typography>
            )}
          </Stack>
        </Stack>
      </BoxWithTitle>

      <form onSubmit={handleSave}>
        <BoxWithTitle title={"editar perfil"} isLoading={isLoading}>
          <Stack direction={"column"} gap={2} sx={{ p: 2 }}>
            <Typography variant={"subtitle2"} fontWeight={"bold"}>dados básicos</Typography>
            <Stack direction={"row"} gap={2} flexWrap={"wrap"}>
              <TextField
                label={"nome"}
                size={"small"}
                required
                value={form.firstName}
                onChange={setField("firstName")}
                slotProps={{ htmlInput: { maxLength: 50 } }}
              />
              <TextField
                label={"sobrenome"}
                size={"small"}
                value={form.lastName}
                onChange={setField("lastName")}
                slotProps={{ htmlInput: { maxLength: 50 } }}
              />
              <TextField
                label={"data de nascimento"}
                type={"date"}
                size={"small"}
                value={form.birthDate}
                onChange={setField("birthDate")}
                slotProps={{ inputLabel: { shrink: true } }}
              />
            </Stack>
            <TextField
              label={"e-mail"}
              size={"small"}
              value={profile?.email ?? ""}
              disabled
              helperText={"O e-mail vem da sua conta Google e não pode ser alterado."}
            />

            <Typography variant={"subtitle2"} fontWeight={"bold"}>sobre você</Typography>
            {profileFields.map(({ key, label, multiline }) => (
              <TextField
                key={key}
                label={label.replace(/:$/, "")}
                size={"small"}
                fullWidth
                multiline={multiline}
                minRows={multiline ? 3 : undefined}
                value={form[key]}
                onChange={setField(key)}
                slotProps={{ htmlInput: { maxLength: 2048 } }}
              />
            ))}

            {saveError && (
              <Typography variant={"body2"} color={"error"}>
                Não foi possível salvar. Confira os campos e tente novamente.
              </Typography>
            )}
            <Stack direction={"row"} gap={1}>
              <Button type={"submit"} loading={isSaving}>salvar</Button>
              <Button onClick={() => navigate(`/profile/${userId}`)}>cancelar</Button>
            </Stack>
          </Stack>
        </BoxWithTitle>
      </form>
    </Stack>
  );
};
