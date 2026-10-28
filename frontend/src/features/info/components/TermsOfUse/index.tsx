import { Stack, Typography } from "@mui/material";
import React from "react";
import { List, Section } from "features/info/components/InfoText";

// Keep this text true to what the app does: update it in the same change as
// anything that stores new data or talks to a new service.

// The terms of use and privacy of Yorgut.
export const TermsOfUse: React.FC = () => (
  <Stack direction={"column"} gap={2}>
    <Typography variant={"body2"}>
      O Yorgut é um projeto pessoal, sem fins lucrativos, mantido por uma única pessoa, para um
      grupo de amigos. Ao usar o Yorgut, você concorda com estes termos.
    </Typography>

    <Section title={"1. Dados que o Yorgut guarda"}>
      <Typography variant={"body2"} gutterBottom>
        O Yorgut guarda somente os dados abaixo. Nenhum outro.
      </Typography>
      <List
        items={[
          "Da sua conta Google, no login: nome, sobrenome, e-mail, foto do perfil e o identificador da conta.",
          "O que você cria: as informações do seu perfil, as fotos que você envia, recados, depoimentos, comunidades, tópicos e respostas no fórum, enquetes.",
          "O que você faz no Yorgut: amizades e pedidos de amizade, comunidades de que participa, convites para comunidades que você envia e recebe, votos em enquetes, avaliações (confiável, legal, sexy) e estrelas de fã que você dá, os perfis que você visita (um registro por perfil e por dia, só enquanto \"registrar minhas visitas\" estiver ligado em Minha conta) e a última vez que você usou o Yorgut.",
          "Suas escolhas em Minha conta: o que quem não é seu amigo vê e a aparência (tema) do Yorgut. O seu navegador também guarda uma cópia do último tema usado, para abrir já com ele.",
          "Registros técnicos do servidor: endereço IP e endereço de cada página pedida, usados somente para encontrar e corrigir erros.",
        ]}
      />
    </Section>

    <Section title={"2. O que o Yorgut nunca faz com seus dados"}>
      <List
        items={[
          "Não vende, não aluga e não compartilha seus dados com nenhuma pessoa, empresa ou plataforma.",
          "Não usa seus dados para treinar inteligência artificial nem qualquer outro modelo, e não permite que ninguém use.",
          "Não usa seus dados para publicidade e não mostra anúncios.",
          "Não usa ferramentas de análise ou rastreamento (como Google Analytics) e não acompanha o que você faz fora do Yorgut.",
        ]}
      />
    </Section>

    <Section title={"3. Quem vê o quê"}>
      <List
        items={[
          "O que você publica (perfil, fotos, recados, depoimentos aprovados, comunidades, fórum e enquetes) pode ser visto por todas as pessoas com conta no Yorgut. O fórum e as enquetes de uma comunidade privada só podem ser vistos pelos membros dela.",
          "Quando você visita um perfil, a pessoa dona do perfil vê o seu nome na lista de visitantes recentes. Para não aparecer, desligue \"registrar minhas visitas\" em Minha conta: suas visitas deixam de ser registradas e seu nome some de todas as listas de visitantes, inclusive das visitas antigas (e você também deixa de ver quem visitou o seu perfil).",
          "Quando você é fã de alguém, o seu nome aparece na lista de fãs dessa pessoa.",
          "Quando você convida um amigo para uma comunidade, ele vê o seu nome no convite.",
          "Suas avaliações não são mostradas a ninguém: os perfis mostram somente as médias.",
          "Seu voto numa enquete não é mostrado a ninguém: as enquetes mostram somente o total de cada opção.",
          "O administrador do Yorgut tem acesso técnico a todos os dados guardados e só os acessa para manter o Yorgut funcionando.",
        ]}
      />
    </Section>

    <Section title={"4. Serviços de terceiros"}>
      <Typography variant={"body2"} gutterBottom>
        O Yorgut não envia seus dados para outras plataformas. Os únicos serviços externos envolvidos
        são estes:
      </Typography>
      <List
        items={[
          "Google, no login: o Google sabe que você entrou no Yorgut, como em qualquer site com login do Google. O Yorgut nunca recebe a sua senha e não acessa nada da sua conta Google além dos dados do item 1.",
          "Google, na foto do perfil: enquanto você não enviar uma foto própria, a foto da sua conta Google continua sendo carregada dos servidores do Google.",
          "Hospedagem: o Yorgut e as fotos enviadas ficam em servidores alugados, que apenas guardam os dados para o Yorgut funcionar.",
        ]}
      />
    </Section>

    <Section title={"5. Como apagar seus dados"}>
      <Typography variant={"body2"}>
        Você pode excluir a sua conta a qualquer momento, em Minha conta. Tudo o que está no item 1 é
        apagado para sempre, inclusive os tópicos que você criou no fórum, com as respostas
        deles, com exceção dos registros técnicos do servidor, que ficam com a hospedagem pelo tempo
        que ela os mantém. As comunidades de que você é dono(a) passam para o membro mais antigo de
        cada uma; as que não têm outros membros são apagadas.
      </Typography>
    </Section>

    <Section title={"6. Regras"}>
      <List
        items={[
          "Não publique nada ilegal, ofensivo ou que exponha outra pessoa sem o consentimento dela.",
          "O administrador pode remover qualquer conteúdo ou conta, a qualquer momento, sem aviso.",
        ]}
      />
    </Section>

    <Section title={"7. Sem garantias"}>
      <Typography variant={"body2"}>
        O Yorgut é oferecido como está, sem garantia de funcionamento, de disponibilidade ou de que
        os dados não serão perdidos. Ele pode sair do ar ou acabar a qualquer momento. Guarde uma
        cópia do que for importante para você.
      </Typography>
    </Section>

    <Section title={"8. Mudanças nestes termos"}>
      <Typography variant={"body2"}>
        Estes termos podem mudar. A versão em vigor é sempre a desta página. Última atualização: 10
        de outubro de 2026.
      </Typography>
    </Section>
  </Stack>
);
