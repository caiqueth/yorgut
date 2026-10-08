import { Box, Stack, Typography } from "@mui/material";
import { useMe } from "api/me/queries";
import {
  useApproveTestimonial,
  useCreateTestimonial,
  useDeleteTestimonial,
  useTestimonials,
} from "api/testimonial/queries";
import { ITestimonial } from "api/testimonial/types";
import { Button } from "components/atoms/Button";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { TextCompositor } from "components/molecules/TextCompositor";
import { dayjsFromDate } from "helpers/formatting/date";
import React, { useState } from "react";
import { useParams } from "react-router-dom";

const testimonialTitle = (t: ITestimonial) =>
  `${t.authorName} em ${dayjsFromDate(t.createdAt, "D [de] MMMM [de] YYYY")}: `;

interface IProps { }
export const UserTestimonialsList: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  const { data: me } = useMe();
  const { data, isLoading } = useTestimonials(userId);
  const { mutate: createTestimonial, isPending: isSending, isSuccess: wasSent } = useCreateTestimonial();
  const { mutate: approveTestimonial, isPending: isApproving } = useApproveTestimonial(userId);
  const { mutate: deleteTestimonial, isPending: isDeleting } = useDeleteTestimonial(userId);

  const [testimonialContent, setTestimonialContent] = useState("");
  const isOwnProfile = me?.id === userId;

  const handleSubmit = () => {
    if (!testimonialContent.trim()) return;
    createTestimonial({
      content: testimonialContent,
      userRecipientId: userId,
    });
  };

  const pendingItems = (data?.pending ?? []).map((t) => (
    <Stack key={t.id} direction={"column"} gap={1}>
      <PicAndTextHorizontal
        imageSrc={t.authorProfilePicture}
        title={testimonialTitle(t)}
        text={t.content}
        url={`/profile/${t.authorId}`}
      />
      <Stack direction={"row"}>
        <Button onClick={() => approveTestimonial(t.id)} loading={isApproving}>aceitar</Button>
        <Button onClick={() => deleteTestimonial(t.id)} loading={isDeleting}>recusar</Button>
      </Stack>
    </Stack>
  ));

  const items = (data?.testimonials ?? []).map((t) => (
    <Stack key={t.id} direction={"column"} gap={1}>
      <PicAndTextHorizontal
        imageSrc={t.authorProfilePicture}
        title={testimonialTitle(t)}
        text={t.content}
        url={`/profile/${t.authorId}`}
      />
      {(isOwnProfile || t.authorId === me?.id) && (
        <Box>
          <Button onClick={() => deleteTestimonial(t.id)} loading={isDeleting}>remover</Button>
        </Box>
      )}
    </Stack>
  ));

  return (
    <Stack direction={"column"} gap={1}>
      {data?.canWrite && (
        <Box>
          <TextCompositor
            content={testimonialContent}
            onChangeContent={setTestimonialContent}
            header={"Escrever um depoimento"}
            onSubmit={handleSubmit}
            isPending={isSending}
            maxLength={2048}
          />
          {wasSent && (
            <Typography variant={"body2"}>
              Depoimento enviado! Ele aparecerá aqui depois de ser aceito.
            </Typography>
          )}
        </Box>
      )}
      {isOwnProfile && pendingItems.length > 0 && (
        <BoxWithTitle title={"depoimentos aguardando aprovação"}>
          <StripedGrid items={pendingItems} elementsPerRow={1} />
        </BoxWithTitle>
      )}
      <BoxWithTitle title={"depoimentos recebidos"}>
        {isLoading || items.length > 0
          ? <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
          : (
            <Box sx={{ p: 1 }}>
              <Typography>Não há nenhum depoimento.</Typography>
            </Box>
          )}
      </BoxWithTitle>
    </Stack>
  );
};
