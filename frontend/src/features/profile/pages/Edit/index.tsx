import { useIdParam } from "hooks/useIdParam";
import { Box, Stack, TextField, Typography } from "@mui/material";
import { useMe, useMeProfile, useUpdateMe, useUpdateProfilePicture } from "api/me/queries";
import { UpdateMeRequest } from "api/generated/model";
import { Button } from "components/Button";
import { BoxWithTitle } from "components/BoxWithTitle";
import { FormattingButtons } from "components/FormattingDialogs";
import { PictureUploadBox } from "components/PictureUploadBox";
import React, { useEffect, useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { profileFields } from "features/profile/profileFields";

const emptyForm: UpdateMeRequest = {
  firstName: "",
  lastName: "",
  birthDate: "",
  ...Object.fromEntries(profileFields.map(({ key }) => [key, ""])),
} as UpdateMeRequest;

// The API sends birth dates as "1990-03-15T00:00:00-03:00". Users created on
// first login get Go's zero time ("0001-01-01..."), shown as empty.
const toDateInput = (birthDate?: string | null) =>
  birthDate && !birthDate.startsWith("0001-") ? birthDate.slice(0, 10) : "";

export const EditProfile: React.FC = () => {
  const userId = useIdParam("userId");
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

  if (me && me.id !== userId) {
    return <Navigate to={`/profile/${userId}`} replace />;
  }

  const setField = (key: keyof UpdateMeRequest) =>
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
      setForm((current) => ({ ...current, [key]: e.target.value }));

  const handleSave = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    updateMe(form, { onSuccess: () => navigate(`/profile/${userId}`) });
  };

  return (
    <Stack direction={"column"} gap={2}>
      <PictureUploadBox
        title={"foto do perfil"}
        noun={"foto"}
        currentPicture={profile?.profilePicture}
        isUploading={isUploading}
        uploadError={uploadError}
        onUpload={(picture, onSent) => updatePicture(picture, { onSuccess: onSent })}
        onPick={resetUpload}
      />

      <form onSubmit={handleSave}>
        <BoxWithTitle title={"editar perfil"} showProgress={isLoading}>
          <Stack direction={"column"} gap={2} sx={{ p: 2 }}>
            <Typography variant={"subtitle2"} fontWeight={"bold"}>dados básicos</Typography>
            <Stack direction={"row"} gap={2} flexWrap={"wrap"}>
              <TextField
                label={"nome"}
                required
                value={form.firstName}
                onChange={setField("firstName")}
                slotProps={{ htmlInput: { maxLength: 50 } }}
              />
              <TextField
                label={"sobrenome"}
                value={form.lastName}
                onChange={setField("lastName")}
                slotProps={{ htmlInput: { maxLength: 50 } }}
              />
              <TextField
                label={"data de nascimento"}
                type={"date"}
                value={form.birthDate}
                onChange={setField("birthDate")}
                slotProps={{ inputLabel: { shrink: true } }}
              />
            </Stack>
            <TextField
              label={"e-mail"}
              value={me?.email ?? ""}
              disabled
              helperText={"O e-mail vem da sua conta Google e não pode ser alterado."}
            />

            <Typography variant={"subtitle2"} fontWeight={"bold"}>sobre você</Typography>
            {profileFields.map(({ key, label, multiline, formatted }) => (
              <Box key={key}>
                <TextField
                  label={label.replace(/:$/, "")}
                  fullWidth
                  multiline={multiline}
                  minRows={multiline ? 3 : undefined}
                  value={form[key]}
                  onChange={setField(key)}
                  slotProps={{ htmlInput: { maxLength: 2048 } }}
                />
                {formatted && <FormattingButtons text={form[key]} />}
              </Box>
            ))}

            {saveError && (
              <Typography variant={"body2"} color={"error"}>
                {saveError.message}
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
