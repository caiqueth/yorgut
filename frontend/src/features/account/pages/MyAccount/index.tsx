import { Box, Checkbox, FormControlLabel, Skeleton, Stack, Switch, Typography } from "@mui/material";
import { PrivacySettings } from "api/generated/model";
import { useMe, useMeProfile, useUpdateMyPrivacy } from "api/me/queries";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Button } from "components/Button";
import { InlineSkeleton } from "components/InlineSkeleton";
import { ProfileTable } from "components/ProfileTable";
import { dayjsFromDate } from "helpers/formatting/date";
import { useLogout } from "hooks/useLogout";
import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ProfileFieldKey, profileFields } from "features/profile/profileFields";
import { ThemePicker } from "features/account/components/ThemePicker";

// The parts of a profile a user may keep to friends.
const sections: { key: "publicAlbum" | "publicCommunities" | "publicScraps" | "publicTestimonials" | "publicBirthday"; label: string }[] = [
  { key: "publicAlbum", label: "álbum de fotos" },
  { key: "publicCommunities", label: "comunidades" },
  { key: "publicScraps", label: "recados" },
  { key: "publicTestimonials", label: "depoimentos" },
  { key: "publicBirthday", label: "aniversário (dia e mês)" },
];

// One column on phones, up to three on wide screens.
const checkboxGridSx = {
  display: "grid",
  gridTemplateColumns: { xs: "1fr", sm: "1fr 1fr", md: "repeat(3, 1fr)" },
};

// The logged user's account (read-only: Google manages it), the app's look and
// privacy settings.
export const MyAccount: React.FC = () => {
  const { data: me } = useMe();
  const { data: profile } = useMeProfile();
  const { mutate: updatePrivacy, isPending: isSaving, error: saveError } = useUpdateMyPrivacy();
  const logout = useLogout();
  const navigate = useNavigate();

  // Unsaved changes, on top of the saved settings. Each box saves its own.
  const [changes, setChanges] = useState<Partial<PrivacySettings>>({});
  const settings = me ? { ...me.privacy, ...changes } : null;
  const change = (change: Partial<PrivacySettings>) => setChanges((prev) => ({ ...prev, ...change }));
  const toggleField = (key: ProfileFieldKey, visible: boolean) => {
    if (!settings) return;
    const hidden = settings.hiddenProfileFields.filter((k) => k !== key);
    change({ hiddenProfileFields: visible ? hidden : [...hidden, key] });
  };

  const visitsKeys: (keyof PrivacySettings)[] = ["recordVisits"];
  const visibilityKeys: (keyof PrivacySettings)[] = [...sections.map((s) => s.key), "hiddenProfileFields"];
  const changed = (keys: (keyof PrivacySettings)[]) =>
    !!me && keys.some((key) => key in changes && !sameValue(changes[key], me.privacy[key]));
  // Saves the changes to keys; the other box keeps its unsaved ones.
  const save = (keys: (keyof PrivacySettings)[]) => {
    if (!me) return;
    const saved = Object.fromEntries(keys.filter((key) => key in changes).map((key) => [key, changes[key]]));
    updatePrivacy({ ...me.privacy, ...saved }, {
      onSuccess: () =>
        setChanges((prev) => Object.fromEntries(Object.entries(prev).filter(([key]) => !keys.includes(key as keyof PrivacySettings)))),
    });
  };

  // The labels and the e-mail note are fixed; only the values wait for the account.
  const loading = <Skeleton variant={"text"} width={"60%"} />;
  const account = [
    { label: "nome:", value: me ? `${me.firstName} ${me.lastName}` : loading },
    {
      label: "e-mail:",
      value: (
        <Stack direction={"column"}>
          {me ? <span>{me.email}</span> : loading}
          <Typography variant={"caption"} color={"text.secondary"}>
            Vem da sua conta Google e não pode ser alterado.
          </Typography>
        </Stack>
      ),
    },
    { label: "membro desde:", value: me ? dayjsFromDate(me.createdAt, "D [de] MMMM [de] YYYY") : loading },
  ];

  return (
    <Stack direction={"column"} gap={2} sx={{ width: "100%" }}>
      <BoxWithTitle title={"minha conta"}>
        <Stack direction={{ xs: "column", sm: "row" }} gap={2} alignItems={"flex-start"} sx={{ p: 2 }}>
          {!profile
            ? <Skeleton variant={"rectangular"} width={96} height={96} />
            : profile.profilePicture
            ? <img src={profile.profilePicture} height={96} alt={"foto do perfil"} />
            // No picture: keep its space, so the table doesn't jump left.
            : <Box sx={{ width: 96, height: 96 }} />}
          <Stack direction={"column"} gap={1} sx={{ flex: 1, width: "100%" }}>
            <ProfileTable data={account} />
            <Typography variant={"body2"}>
              Nome, foto e o resto do perfil se mudam em{" "}
              {me ? <Link to={`/profile/${me.id}/edit`}>editar perfil</Link> : <InlineSkeleton width={"11ch"} />}.
            </Typography>
            <Box>
              <Button onClick={logout}>sair</Button>
            </Box>
          </Stack>
        </Stack>
      </BoxWithTitle>

      <ThemePicker />

      <BoxWithTitle title={"privacidade"} showProgress={isSaving}>
        <Stack direction={"column"} gap={1} sx={{ p: 2 }}>
          <FormControlLabel
            control={
              <Switch
                checked={settings?.recordVisits ?? true}
                disabled={!settings || isSaving}
                onChange={(e) => change({ recordVisits: e.target.checked })}
              />
            }
            label={"registrar minhas visitas"}
          />
          <Typography variant={"body2"}>
            Ligado: quando você visita um perfil, a pessoa dona dele vê o seu nome nos visitantes
            recentes, e você vê quem visitou o seu.
          </Typography>
          <Typography variant={"body2"}>
            Desligado: suas visitas não são registradas, seu nome some da lista de visitantes de todo
            mundo (inclusive das visitas antigas), e você também não vê quem visitou o seu perfil.
          </Typography>
          <Box>
            <Button disabled={!changed(visitsKeys)} loading={isSaving} onClick={() => save(visitsKeys)}>salvar</Button>
          </Box>
          {saveError && (
            <Typography variant={"body2"} color={"error"}>
              {saveError.message}
            </Typography>
          )}
        </Stack>
      </BoxWithTitle>

      <BoxWithTitle title={"o que quem não é seu amigo vê"} showProgress={isSaving}>
        <Stack direction={"column"} gap={2} sx={{ p: 2 }}>
          <Typography variant={"body2"}>
            Seus amigos sempre veem tudo. Marque o que as outras pessoas também podem ver no seu perfil.
          </Typography>
          <Box sx={checkboxGridSx}>
            {sections.map(({ key, label }) => (
              <FormControlLabel
                key={key}
                control={
                  <Checkbox
                    size={"small"}
                    checked={settings?.[key] ?? false}
                    disabled={!settings || isSaving}
                    onChange={(e) => change({ [key]: e.target.checked })}
                  />
                }
                label={label}
              />
            ))}
          </Box>
          <Typography variant={"body2"}>
            Com os recados desmarcados, quem não é seu amigo também não pode deixar recados para você.
          </Typography>
          <Typography variant={"subtitle2"} fontWeight={"bold"}>campos do perfil</Typography>
          <Box sx={checkboxGridSx}>
            {profileFields.map(({ key, label }) => (
              <FormControlLabel
                key={key}
                control={
                  <Checkbox
                    size={"small"}
                    checked={!!settings && !settings.hiddenProfileFields.includes(key)}
                    disabled={!settings || isSaving}
                    onChange={(e) => toggleField(key, e.target.checked)}
                  />
                }
                label={label.replace(/:$/, "")}
              />
            ))}
          </Box>
          <Stack direction={"row"} gap={1} flexWrap={"wrap"}>
            <Button disabled={!changed(visibilityKeys)} loading={isSaving} onClick={() => save(visibilityKeys)}>
              salvar
            </Button>
            <Button
              disabled={!me}
              onClick={() => me && navigate(`/profile/${me.id}`, { state: { preview: true } })}
            >
              ver meu perfil como visitante
            </Button>
          </Stack>
        </Stack>
      </BoxWithTitle>
      <BoxWithTitle title={"excluir conta"}>
        <Stack direction={"column"} gap={1} sx={{ p: 2 }} alignItems={"flex-start"}>
          <Typography variant={"body2"}>
            Apaga a sua conta e tudo o que você fez no Yorgut, para sempre.
          </Typography>
          <Button onClick={() => navigate("/my-account/delete")}>excluir minha conta</Button>
        </Stack>
      </BoxWithTitle>
    </Stack>
  );
};

// Hidden fields count as unchanged when only their order differs.
const sameValue = (a: unknown, b: unknown) =>
  Array.isArray(a) && Array.isArray(b)
    ? a.length === b.length && a.every((item) => b.includes(item))
    : a === b;
