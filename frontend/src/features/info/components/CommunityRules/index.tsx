import { Stack, Typography } from "@mui/material";
import React from "react";
import { List, Section } from "features/info/components/InfoText";
import { useInstanceInfo } from "api/info/queries";
import { InlineSkeleton } from "components/InlineSkeleton";

// Yorgut's community rules, shown on the help page. The tools they point to
// (deleting scraps, approving testimonials, "só para amigos", undoing a
// friendship) must exist: keep them in sync with the app, as for the terms.
export const CommunityRules: React.FC = () => (
  <Stack direction={"column"} gap={2}>
    <Typography variant={"body2"}>
      O Yorgut é pequeno e feito para amigos. Estas regras valem para tudo o que você escreve e
      envia aqui: perfil, fotos, recados, depoimentos, comunidades, fórum e enquetes. A ideia é
      simples: tratar o Yorgut como a sala de casa, cheia de gente que você conhece.
    </Typography>

    <Section title={"1. Do outro lado da tela tem uma pessoa"}>
      <List
        items={[
          "Antes de enviar, pense: eu diria isso pessoalmente? Discorde de ideias, não ataque pessoas.",
          "Nada de ofensas, ameaças, assédio ou perseguição, nem de insistir com quem já pediu para você parar.",
          "Nada de discurso de ódio ou discriminação por raça, etnia, origem, religião, gênero, orientação sexual, deficiência, idade ou aparência.",
          "Não exponha ninguém: nada de publicar endereço, telefone, documentos, conversas privadas ou fotos de outra pessoa sem a permissão dela.",
        ]}
      />
    </Section>

    <Section title={"2. Presuma boa-fé"}>
      <List
        items={[
          "Texto não tem tom de voz. Na dúvida, pense que a pessoa quis dizer a coisa mais gentil possível, e pergunte antes de se irritar.",
          "Escreva de um jeito claro. Ninguém precisa ser escritor, mas ajuda reler antes de enviar.",
        ]}
      />
    </Section>

    <Section title={"3. Seja você"}>
      <List
        items={[
          "Use seu nome ou o apelido pelo qual seus amigos te conhecem.",
          "Não se passe por outra pessoa e não crie perfis falsos.",
        ]}
      />
    </Section>

    <Section title={"4. O que não pode ser publicado"}>
      <List
        items={[
          "Nada ilegal, e nada que envolva menores de idade de forma sexual, nunca.",
          "Nada de nudez, pornografia ou violência explícita, nem nas fotos nem nos textos.",
          "Respeite o trabalho dos outros: publique fotos e textos que são seus ou que você pode compartilhar, dizendo de quem são.",
          "Nada de golpes, links maliciosos ou boatos que podem fazer mal a alguém.",
        ]}
      />
    </Section>

    <Section title={"5. Sem spam"}>
      <List
        items={[
          "Nada de propaganda, correntes, mensagens em massa ou links de afiliado e de indicação.",
          "Comunidades, tópicos e enquetes devem ser de verdade, com um assunto, e não para divulgar algo.",
        ]}
      />
    </Section>

    <Section title={"6. Comunidades"}>
      <List
        items={[
          "Quem cria a comunidade cuida dela: escreve a descrição, decide quem entra numa comunidade privada e pode encerrar enquetes.",
          "Mantenha os tópicos e as enquetes no assunto da comunidade.",
          "Convide só amigos que você acha que vão gostar. Convite não é propaganda.",
        ]}
      />
    </Section>

    <Section title={"7. Inteligência artificial"}>
      <List
        items={[
          "Tudo bem usar ferramentas como o ChatGPT para traduzir ou revisar um texto seu.",
          "Mas recados e depoimentos são sobre amizades de verdade: escreva com as suas palavras, não com as de uma máquina.",
        ]}
      />
    </Section>

    <Section title={"8. Quando algo dá errado"}>
      <Typography variant={"body2"} gutterBottom>
        Você tem ferramentas para cuidar do seu espaço:
      </Typography>
      <List
        items={[
          "Apague qualquer recado do seu perfil, e os que você escreveu em outros.",
          "Depoimentos só aparecem depois que você aprova. Você pode recusar ou apagar os que recebeu.",
          "Em Minha conta, deixe álbum, recados, depoimentos, comunidades e aniversário visíveis só para amigos.",
          "Desfaça uma amizade quando quiser.",
        ]}
      />
      <Typography variant={"body2"} sx={{ mt: 1 }}>
        Se alguém quebrar estas regras, avise o administrador do Yorgut
        <AdminEmail />. Ele pode apagar o que quebrar as
        regras e, em casos graves ou repetidos, suspender ou excluir a conta, avisando por e-mail
        sempre que possível.
      </Typography>
    </Section>

    <Typography variant={"body2"} color={"text.secondary"}>
      Atualizadas em 15 de outubro de 2026.
    </Typography>
  </Stack>
);

// " pelo e-mail x@y.com", from the API (ADMIN_EMAIL); nothing when the
// instance has none.
const AdminEmail: React.FC = () => {
  const { data, isLoading } = useInstanceInfo();
  if (isLoading) return <> pelo e-mail <InlineSkeleton width={"16ch"} /></>;
  const email = data?.adminEmail;
  if (!email) return null;
  return <> pelo e-mail <a href={`mailto:${email}`}>{email}</a></>;
};
