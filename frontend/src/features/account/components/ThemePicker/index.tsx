import { FormControlLabel, Radio, RadioGroup, Skeleton, Stack, Typography } from "@mui/material";
import { useMe, useUpdateMyTheme } from "api/me/queries";
import { BoxWithTitle } from "components/BoxWithTitle";
import React from "react";
import { isThemeKey, themeOptions } from "theme";

// The look of the app. Saved on pick: the theme shows at once, so there's
// nothing to confirm.
export const ThemePicker: React.FC = () => {
  const { data: me } = useMe();
  const { mutate: updateTheme, isPending, error } = useUpdateMyTheme();

  return (
    <BoxWithTitle title={"aparência"} showProgress={isPending}>
      <Stack direction={"column"} gap={1} sx={{ p: 2 }}>
        <Typography variant={"body2"}>
          Escolha o visual do yorgut. Ele fica salvo na sua conta e vale em qualquer navegador.
        </Typography>
        {!me ? (
          themeOptions.map(({ key }) => <Skeleton key={key} variant={"text"} width={"40%"} />)
        ) : (
          <RadioGroup
            value={me.theme}
            onChange={(e) => isThemeKey(e.target.value) && updateTheme(e.target.value)}
          >
            {themeOptions.map(({ key, label, description }) => (
              <FormControlLabel
                key={key}
                value={key}
                control={<Radio size={"small"} />}
                label={
                  <Stack direction={"column"}>
                    <Typography variant={"body2"} fontWeight={"bold"}>{label}</Typography>
                    <Typography variant={"caption"} color={"text.secondary"}>{description}</Typography>
                  </Stack>
                }
                sx={{ alignItems: "flex-start", mb: 1, "& .MuiRadio-root": { pt: 0.25 } }}
              />
            ))}
          </RadioGroup>
        )}
        {error && (
          <Typography variant={"body2"} color={"error"}>
            {error.message}
          </Typography>
        )}
      </Stack>
    </BoxWithTitle>
  );
};
