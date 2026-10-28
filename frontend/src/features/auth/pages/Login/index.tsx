import { Box, CircularProgress, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";
import { CredentialResponse, GoogleLogin } from "@react-oauth/google";
import { Link, useNavigate } from "react-router-dom";
import { useLoginWithGoogle } from "api/auth/queries";
import whoDoYouKnow from "assets/images/who_do_you_know.jpg";
import { Disclaimer } from "components/Disclaimer";

export const Login = () => {
  const navigate = useNavigate();
  const loginWithGoogle = useLoginWithGoogle();

  const onSuccess = (credentialResponse: CredentialResponse) => {
    if (!credentialResponse.credential) {
      onError();
      return;
    }
    loginWithGoogle.mutate(
      { credential: credentialResponse.credential },
      { onSuccess: () => navigate("/") },
    );
  };

  const onError = () => {
    console.log("Login Failed");
  };

  return (
    <Grid container spacing={4}>
      <Grid size={{ xs: 12, md: 8 }}>
        <Paper sx={{ p: 4, backgroundColor: "white", height: "100%", boxSizing: "border-box" }}>
          <Stack direction={"column"} alignItems={"center"} gap={1}>
            <Typography variant={"body2"}>
              Yorgut é uma comunidade online que conecta pessoas através de uma
              rede de amigos confiáveis.
            </Typography>
            <Typography variant={"body2"}>
              Estamos comprometidos em oferecer um espaço de encontro online
              onde as pessoas possam socializar, fazer novas amizades e
              encontrar outras que compartilhem seus interesses.
            </Typography>
            <img src={whoDoYouKnow} height={136} width={235} alt={"quem você conhece?"} />
            <Typography variant={"body2"}>
              Junte-se ao Yorgut para expandir a circunferência do seu círculo
              social.
            </Typography>
            <Box sx={{ mt: 2 }}>
              <Disclaimer variant={"full"} muted />
            </Box>
          </Stack>
        </Paper>
      </Grid>
      {/* Grid items stretch to the row's height (set by the left box); the login
          box fills its item, and its white area takes the rest, content centred. */}
      {/* On phones the login box comes first, above the introduction. */}
      <Grid size={{ xs: 12, md: 4 }} sx={{ order: { xs: -1, md: 0 } }}>
        <Paper sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
          <Box sx={{ p: 1 }}>
            <Typography>login</Typography>
          </Box>
          <Paper
            sx={{
              m: 1,
              p: 1,
              backgroundColor: "white",
              flex: 1,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <Stack direction={"column"} alignItems={"center"} gap={3}>
              <Typography variant={"body2"}>
                Faça login no Yorgut com a sua conta Google.
              </Typography>
              {loginWithGoogle.isPending
                ? (
                  <Stack direction={"row"} alignItems={"center"} gap={1}>
                    <CircularProgress size={20} />
                    <Typography variant={"body2"}>Entrando...</Typography>
                  </Stack>
                )
                : <GoogleLogin onSuccess={onSuccess} onError={onError} />}
              {loginWithGoogle.error && (
                <Typography variant={"body2"} color={"error"}>
                  {loginWithGoogle.error.message}
                </Typography>
              )}
              {/* What "Sign in with Google" gives the API (AuthService.LoginWithGoogle
                  reads only these claims of the ID token; there is no access token,
                  so nothing else in the account can be reached). Keep it in sync. */}
              <Typography variant={"body2"} color={"text.secondary"} sx={{ textAlign: "justify", hyphens: "auto" }} lang={"pt-BR"}>
                O Google só informa ao Yorgut o seu nome, e-mail e foto de perfil (e um código da sua
                conta, para reconhecer você no próximo login). O Yorgut não vê a sua senha e não tem acesso
                ao seu Gmail, contatos, Drive ou qualquer outra coisa da sua conta Google.{" "}
                <Link to={"/about"}>Veja tudo o que guardamos.</Link>
              </Typography>
              <Typography variant={"body2"}>
                É necessário ter uma conta Google para acessar o Yorgut.
              </Typography>
            </Stack>
          </Paper>
        </Paper>
      </Grid>
    </Grid>
  );
};
