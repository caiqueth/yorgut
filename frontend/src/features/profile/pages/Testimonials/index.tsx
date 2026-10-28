import { useIdParam } from "hooks/useIdParam";
import { Box, Stack, Typography } from "@mui/material";
import { useMe } from "api/me/queries";
import { useUserProfileSection } from "api/user/queries";
import { PrivateSection } from "components/PrivateSection";
import { useProfilePreview } from "contexts/ProfilePreviewContext";
import {
  useApproveTestimonial,
  useCreateTestimonial,
  useDeleteTestimonial,
  useTestimonials,
} from "api/testimonial/queries";
import { TestimonialItem } from "api/generated/model";
import { Button } from "components/Button";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { TextCompositor } from "components/TextCompositor";
import { ListPagination } from "components/ListPagination";
import { TESTIMONIALS_PAGE_SIZE } from "helpers/paging";
import { usePage } from "hooks/usePage";
import { dayjsFromDate } from "helpers/formatting/date";
import React, { useState } from "react";

// After the author's name (the link): " em 12 de outubro de 2026:".
const testimonialDate = (t: TestimonialItem) =>
  ` em ${dayjsFromDate(t.createdAt, "D [de] MMMM [de] YYYY")}:`;

export const UserTestimonialsList: React.FC = () => {
  const userId = useIdParam("userId");

  const { data: me } = useMe();
  const { previewing } = useProfilePreview();
  const section = useUserProfileSection(userId, "testimonials");
  const [page, setPage] = usePage(userId);
  const { data, isLoading: isLoadingTestimonials, isFetching } = useTestimonials(userId, page, TESTIMONIALS_PAGE_SIZE, section.canSee);
  // The request waits for the profile, which says whether it may be made.
  const isLoading = section.isLoading || isLoadingTestimonials;
  const { mutate: createTestimonial, isPending: isSending, isSuccess: wasSent } = useCreateTestimonial();
  const { mutate: approveTestimonial, isPending: isApproving, variables: approvingId } = useApproveTestimonial(userId);
  const { mutate: deleteTestimonial, isPending: isDeleting, variables: deletingId } = useDeleteTestimonial(userId);

  const [testimonialContent, setTestimonialContent] = useState("");
  // Not in the preview, which shows a stranger's view.
  const isOwnProfile = me?.id === userId && !previewing;

  const handleSubmit = () => {
    if (!testimonialContent.trim()) return;
    createTestimonial(
      {
        content: testimonialContent,
        userRecipientId: userId,
      },
      { onSuccess: () => setTestimonialContent("") },
    );
  };

  const pendingItems = (data?.pending ?? []).map((t) => (
    <Stack key={t.id} direction={"column"} gap={1}>
      <PicAndTextHorizontal
        imageSrc={t.authorProfilePicture}
        title={t.authorName}
        titleSuffix={testimonialDate(t)}
        text={t.content}
        url={`/profile/${t.authorId}`}
      />
      <Stack direction={"row"}>
        <Button onClick={() => approveTestimonial(t.id)} loading={isApproving && approvingId === t.id}>aceitar</Button>
        <Button onClick={() => deleteTestimonial(t.id)} loading={isDeleting && deletingId === t.id}>recusar</Button>
      </Stack>
    </Stack>
  ));

  const items = (data?.testimonials ?? []).map((t) => (
    <Stack key={t.id} direction={"column"} gap={1}>
      <PicAndTextHorizontal
        imageSrc={t.authorProfilePicture}
        title={t.authorName}
        titleSuffix={testimonialDate(t)}
        text={t.content}
        url={`/profile/${t.authorId}`}
      />
      {(isOwnProfile || t.authorId === me?.id) && (
        <Box>
          <Button onClick={() => deleteTestimonial(t.id)} loading={isDeleting && deletingId === t.id}>remover</Button>
        </Box>
      )}
    </Stack>
  ));

  if (!section.isLoading && !section.canSee) {
    return (
      <BoxWithTitle title={"depoimentos"}>
        <PrivateSection firstName={section.firstName} what={"os depoimentos"} />
      </BoxWithTitle>
    );
  }

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
        <BoxWithTitle title={"depoimentos aguardando aprovação"} showProgress={isFetching && !isLoading}>
          <StripedGrid items={pendingItems} elementsPerRow={1} />
        </BoxWithTitle>
      )}
      <BoxWithTitle title={"depoimentos recebidos"} showProgress={isFetching && !isLoading}>
        {isLoading || items.length > 0
          ? <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<PicAndTextHorizontalSkeleton lineWidths={[240, 200]} />} />
          : (
            <Box sx={{ p: 1 }}>
              <Typography>Não há nenhum depoimento.</Typography>
            </Box>
          )}
        <ListPagination page={page} pageSize={TESTIMONIALS_PAGE_SIZE} count={data?.count} onChange={setPage} />
      </BoxWithTitle>
    </Stack>
  );
};
