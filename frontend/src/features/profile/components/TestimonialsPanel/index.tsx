import { Box, Pagination, Stack, Typography } from "@mui/material";
import { useTestimonials } from "api/testimonial/queries";
import { PicAndTextHorizontal } from "components/PicAndTextHorizontal";
import { BoxWithTitle } from "components/BoxWithTitle";
import { StripedGrid } from "components/StripedGrid";
import { dayjsFromDate } from "helpers/formatting/date";
import React, { useState } from "react";
import { Link } from "react-router-dom";

const PAGE_SIZE = 5;

interface IProps {
  userId: number;
}
export const TestimonialsPanel: React.FC<IProps> = ({ userId }) => {
  const [page, setPage] = useState(1);
  const { data, isLoading, isFetching } = useTestimonials(userId, page, PAGE_SIZE);

  // Start over when navigating to another profile.
  React.useEffect(() => setPage(1), [userId]);

  const count = data?.count ?? 0;
  const pageCount = Math.ceil(count / PAGE_SIZE);

  const items = (data?.testimonials ?? []).map((t) => (
    <PicAndTextHorizontal
      key={t.id}
      imageSrc={t.authorProfilePicture}
      title={`${t.authorName} em ${dayjsFromDate(t.createdAt, "D [de] MMMM [de] YYYY")}: `}
      text={t.content}
      url={`/profile/${t.authorId}`}
    />
  ));

  return (
    <BoxWithTitle title={`depoimentos (${count})`} isLoading={isFetching && !isLoading}>
      {isLoading || items.length > 0
        ? <StripedGrid items={items} elementsPerRow={1} isLoading={isLoading} />
        : (
          <Box sx={{ p: 1 }}>
            <Typography variant={"body2"}>Nenhum depoimento ainda.</Typography>
          </Box>
        )}
      <Stack direction={"row"} justifyContent={"space-between"} alignItems={"center"} sx={{ p: 1 }}>
        {pageCount > 1
          ? (
            <Pagination
              size={"small"}
              count={pageCount}
              page={page}
              onChange={(_e, value) => setPage(value)}
            />
          )
          : <span />}
        <Typography variant={"body2"}>
          <Link to={`/profile/${userId}/testimonials`}>ver todos</Link>
        </Typography>
      </Stack>
    </BoxWithTitle>
  );
};
