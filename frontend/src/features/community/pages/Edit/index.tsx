import { useIdParam } from "hooks/useIdParam";
import { Box, Checkbox, FormControlLabel, Stack, TextField, Typography } from "@mui/material";
import {
  useCommunity,
  useCommunityProfile,
  useUpdateCommunity,
  useUpdateCommunityPicture,
} from "api/community/queries";
import { CommunityCategory, UpdateCommunityRequest } from "api/generated/model";
import { useMe } from "api/me/queries";
import { Button } from "components/Button";
import { BoxWithTitle } from "components/BoxWithTitle";
import { FormattingButtons } from "components/FormattingDialogs";
import { CommunityCategorySelect } from "components/CommunityCategorySelect";
import { PictureUploadBox } from "components/PictureUploadBox";
import React, { useEffect, useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";

// The form keeps "" as the category until the community loads.
type IForm = Omit<UpdateCommunityRequest, "category"> & { category: CommunityCategory | "" };
export const EditCommunity: React.FC = () => {
  const communityId = useIdParam("communityId");
  const navigate = useNavigate();

  const { data: me } = useMe();
  const { data: community, isLoading } = useCommunity(communityId);
  const { data: profile } = useCommunityProfile(communityId);
  const { mutate: updateCommunity, isPending: isSaving, error: saveError } = useUpdateCommunity(communityId);
  const {
    mutate: updatePicture,
    isPending: isUploading,
    error: uploadError,
    reset: resetUpload,
  } = useUpdateCommunityPicture(communityId);

  const [form, setForm] = useState<IForm>({ name: "", category: "", description: "", private: false });

  // Fill the form once the current data arrives.
  useEffect(() => {
    if (!community || !profile) return;
    setForm({
      name: community.name ?? "",
      category: community.category,
      description: profile.description ?? "",
      private: community.private,
    });
  }, [community, profile]);

  // Only the owner can edit; the backend enforces it too.
  if (me && community && me.id !== community.ownerId) {
    return <Navigate to={`/community/${communityId}`} replace />;
  }

  const setField = (key: "name" | "description") =>
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
      setForm((current) => ({ ...current, [key]: e.target.value }));

  const handleSave = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // The select is required, so the browser doesn't submit without a category.
    if (!form.category) return;
    updateCommunity(
      { ...form, category: form.category },
      { onSuccess: () => navigate(`/community/${communityId}`) },
    );
  };

  return (
    <Stack direction={"column"} gap={2}>
      <PictureUploadBox
        title={"imagem da comunidade"}
        noun={"imagem"}
        currentPicture={profile?.profilePicture}
        isUploading={isUploading}
        uploadError={uploadError}
        onUpload={(picture, onSent) => updatePicture(picture, { onSuccess: onSent })}
        onPick={resetUpload}
      />

      <form onSubmit={handleSave}>
        <BoxWithTitle title={"editar comunidade"} showProgress={isLoading}>
          <Stack direction={"column"} gap={2} sx={{ p: 2 }}>
            <TextField
              label={"nome"}
              required
              value={form.name}
              onChange={setField("name")}
              slotProps={{ htmlInput: { maxLength: 100 } }}
            />
            <CommunityCategorySelect
              label={"categoria"}
              value={form.category}
              onChange={(category) => setForm((current) => ({ ...current, category }))}
            />
            <Box>
              <TextField
                label={"descrição"}
                multiline
                minRows={4}
                fullWidth
                value={form.description}
                onChange={setField("description")}
                slotProps={{ htmlInput: { maxLength: 2048 } }}
              />
              <FormattingButtons text={form.description} />
            </Box>
            <Box>
              <FormControlLabel
                control={
                  <Checkbox
                    size={"small"}
                    checked={form.private}
                    onChange={(e) => setForm((current) => ({ ...current, private: e.target.checked }))}
                  />
                }
                label={"comunidade privada: só membros veem o fórum, e você aprova quem entra"}
              />
              {community?.private && !form.private && (
                <Typography variant={"body2"}>
                  Ao tornar a comunidade pública, todos os pedidos pendentes serão aceitos.
                </Typography>
              )}
            </Box>
            {saveError && (
              <Typography variant={"body2"} color={"error"}>
                {saveError.message}
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
