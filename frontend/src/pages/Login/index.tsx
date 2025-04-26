import { Box, Grid2 as Grid, Paper, Stack, Typography } from "@mui/material";
import { CredentialResponse, GoogleLogin } from "@react-oauth/google";
import { useNavigate } from "react-router-dom";
import whoDoYouKnow from "../../assets/images/who_do_you_know.jpg";

const onSuccess = (credentialResponse: CredentialResponse) => {
  const navigate = useNavigate();
  localStorage.setItem("token", JSON.stringify(credentialResponse.credential));
  navigate("/");
};

const onError = () => {
  console.log("Login Failed");
};

export const Login = () => {
  return (
    <Grid container spacing={4}>
      <Grid size={8}>
        <Paper sx={{ p: 4, backgroundColor: "white" }}>
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
            <img src={whoDoYouKnow} height={136} width={235}></img>
            <Typography variant={"body2"}>
              Junte-se ao Yorgut para expandir a circunferência do seu círculo
              social.
            </Typography>
          </Stack>
        </Paper>
      </Grid>
      <Grid size={4}>
        <Paper>
          <Box sx={{ p: 1 }}>
            <Typography>login</Typography>
          </Box>
          <Paper sx={{ m: 1, p: 1, backgroundColor: "white" }}>
            <Stack direction={"column"} alignItems={"center"} gap={3}>
              <Typography variant={"body2"}>
                Faça login no Yorgut com a sua conta Google.
              </Typography>
              <GoogleLogin onSuccess={onSuccess} onError={onError} />
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
