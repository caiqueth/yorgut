import { Box, Stack, TextField, Typography } from "@mui/material";
import { Button } from "components/atoms/Button";
import { PicAndTextHorizontal } from "components/atoms/PicAndTextHorizontal";
import { BoxWithTitle } from "components/molecules/BoxWithTitle";
import { StripedGrid } from "components/molecules/StripedGrid";
import { TextCompositor } from "components/molecules/TextCompositor";
import React, { useState } from "react";
import { useParams } from "react-router-dom";

interface IProps { }
export const UserTestimonialsList: React.FC<IProps> = () => {
  const params = useParams();
  const userId = parseInt(params.userId ?? "");

  const [testimonialContent, setTestimonialContent] = useState("");
  const handleSubmit = () => { };

  const { data: testimonials } = {
    data: {
      testimonials: [
        {
          author: "👷",
          text: " Em construção...",
        },
      ],
    },
  };

  const items = testimonials.testimonials.map((t) => (
    <PicAndTextHorizontal
      title={t.author}
      text={t.text}
      imageSrc={t.author}
      url={t.author}
    />
  ));

  return (
    <Box>
      <Box>
        <TextCompositor
          content={testimonialContent}
          onChangeContent={setTestimonialContent}
          header={"Escrever um depoimento"}
          onSubmit={handleSubmit}
          isPending={true}
          maxLength={2048}
        />
      </Box>
      <BoxWithTitle title={"depoimentos recebidos"}>
        <StripedGrid items={items} elementsPerRow={1} isLoading={true} />
      </BoxWithTitle>
    </Box>
  );
};
