import {
  Box,
  Stack,
  Typography,
} from "@mui/material";
import { useCreateScrap, useDeleteScrap, useScraps } from "api/scrap/queries";
import { useMe } from "api/me/queries";
import { useUserProfileSection } from "api/user/queries";
import { PrivateSection } from "components/PrivateSection";
import { useProfilePreview } from "contexts/ProfilePreviewContext";
import { Button } from "components/Button";
import { PicAndTextHorizontal, PicAndTextHorizontalSkeleton } from "components/PicAndTextHorizontal";

import { BoxWithTitle, Count } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { ListPagination } from "components/ListPagination";
import { SCRAP_PAGE_SIZE } from "helpers/paging";
import { TextCompositor } from "components/TextCompositor";
import { FormattedText } from "components/FormattedText";
import { dayjsFromDate } from "helpers/formatting/date";
import React, { useState } from "react";

interface IProps {
  userId: number;
}
export const Scrapbook: React.FC<IProps> = ({ userId }) => {
  const [page, setPage] = useState(1);
  // A scrapbook kept to friends is neither shown nor written to by others.
  const section = useUserProfileSection(userId, "scraps");
  const { previewing } = useProfilePreview();
  const { data: scraps, isLoading: isLoadingScraps, isFetching } = useScraps(userId, page, SCRAP_PAGE_SIZE, section.canSee);
  // The request waits for the profile, which says whether it may be made.
  const isLoading = section.isLoading || isLoadingScraps;
  const { mutate: sendScrap, isPending } = useCreateScrap();
  const { mutate: removeScrap, isPending: isDeleting, variables: deletingId } = useDeleteScrap(userId);
  const { data: me } = useMe();
  // The scrap whose "apagar" was clicked, waiting for "sim" or "não".
  const [confirmingId, setConfirmingId] = useState<number | null>(null);

  const [scrapContent, setScrapContent] = useState<string>("");

  // Start over when navigating to another scrapbook.
  React.useEffect(() => setPage(1), [userId]);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    sendScrap(
      {
        content: scrapContent,
        userRecipientId: userId,
      },
      {
        // The new scrap is the newest: show it on the first page.
        onSuccess: () => {
          setScrapContent("");
          setPage(1);
        },
      },
    );
  };

  // The scrapbook's owner and the scrap's author may delete it (the API
  // checks the same); nobody in the preview, which shows a stranger's view.
  const canDelete = (authorId: number) => !!me && !previewing && (me.id === userId || me.id === authorId);

  const deleteControls = (scrapId: number) => {
    if (isDeleting && deletingId === scrapId) {
      return <Button size={"small"} loading>apagar</Button>;
    }
    if (confirmingId === scrapId) {
      return (
        <Stack direction={"row"} alignItems={"center"}>
          <Typography variant={"body2"}>apagar?</Typography>
          <Button
            size={"small"}
            onClick={() => removeScrap(scrapId, { onSettled: () => setConfirmingId(null) })}
          >
            sim
          </Button>
          <Button size={"small"} onClick={() => setConfirmingId(null)}>não</Button>
        </Stack>
      );
    }
    return <Button size={"small"} onClick={() => setConfirmingId(scrapId)}>apagar</Button>;
  };

  const items =
    scraps?.scraps.map((s) => (
      <Stack
        key={s.id}
        direction={"row"}
        justifyContent={"space-between"}
        alignItems={"flex-start"}
        gap={1}
      >
        <PicAndTextHorizontal
          imageSrc={s.authorProfilePicture}
          title={s.authorName}
          titleSuffix={` em ${dayjsFromDate(s.createdAt, "D [de] MMMM [de] YYYY [às] HH[h]mm")}:`}
          text={<FormattedText text={s.content} />}
          url={`/profile/${s.authorId}`}
        />
        {canDelete(s.authorId) && <Box sx={{ flexShrink: 0 }}>{deleteControls(s.id)}</Box>}
      </Stack>
    )) ?? [];

  if (!section.isLoading && !section.canSee) {
    return (
      <BoxWithTitle title={"recados"}>
        <PrivateSection firstName={section.firstName} what={"os recados"} />
      </BoxWithTitle>
    );
  }

  return (
    <Stack direction={"column"} gap={1}>
      <TextCompositor
        content={scrapContent}
        header={"Deixar um recado"}
        onChangeContent={setScrapContent}
        onSubmit={handleSubmit}
        isPending={isPending}
        maxLength={2048}
        // Shown in the preview as strangers see it, but it doesn't send.
        canSubmit={!previewing}
        formatted
      />
      <BoxWithTitle title={<>recados <Count value={scraps?.count} /></>} showProgress={isFetching && !isLoading}>
        {isLoading || items.length > 0
          ? <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} skeleton={<PicAndTextHorizontalSkeleton lineWidths={[280, 220]} />} />
          : (
            <Box sx={{ p: 1 }}>
              <Typography>Não há nenhum recado.</Typography>
            </Box>
          )}
        <ListPagination page={page} pageSize={SCRAP_PAGE_SIZE} count={scraps?.count} onChange={setPage} />
      </BoxWithTitle>
    </Stack>
  );
};
