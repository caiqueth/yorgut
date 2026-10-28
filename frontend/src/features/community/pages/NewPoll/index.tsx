import React, { useState } from "react";
import dayjs from "dayjs";
import { Box, Stack, TextField, Typography } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useCreatePoll } from "api/poll/queries";
import { BoxWithTitle } from "components/BoxWithTitle";
import { Button } from "components/Button";
import { MembersOnly } from "features/community/components/MembersOnly";
import { useIdParam } from "hooks/useIdParam";

// The API's limits.
const MIN_OPTIONS = 2;
const MAX_OPTIONS = 10;

export const NewPoll: React.FC = () => {
  const communityId = useIdParam("communityId");
  const navigate = useNavigate();
  const { mutate: createPoll, isPending, error } = useCreatePoll(communityId);
  const [question, setQuestion] = useState("");
  const [description, setDescription] = useState("");
  const [options, setOptions] = useState<string[]>(["", ""]);
  // "YYYY-MM-DD" from the date input; empty keeps the poll open until closed.
  const [closesOn, setClosesOn] = useState("");

  const setOption = (index: number, text: string) =>
    setOptions((current) => current.map((o, i) => (i === index ? text : o)));
  const removeOption = (index: number) => setOptions((current) => current.filter((_, i) => i !== index));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createPoll(
      {
        question,
        description,
        options,
        // Voting ends at the end of the chosen day, in the user's time zone.
        closesAt: closesOn ? dayjs(closesOn).endOf("day").toISOString() : undefined,
      },
      { onSuccess: (response) => navigate(`/community/${communityId}/polls/${response.data!.id}`) },
    );
  };

  return (
    <MembersOnly communityId={communityId} action={"criar enquetes"}>
      <BoxWithTitle title={"nova enquete"}>
        <form onSubmit={handleSubmit}>
          <Stack sx={{ p: 1 }} gap={2}>
            <TextField
              label={"pergunta"}
              required
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              slotProps={{ htmlInput: { maxLength: 200 } }}
            />
            <TextField
              label={"descrição (opcional)"}
              multiline
              minRows={2}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              slotProps={{ htmlInput: { maxLength: 2048 } }}
            />
            <Box>
              <Typography>Opções</Typography>
              <Stack gap={1} sx={{ mt: 1 }}>
                {options.map((o, i) => (
                  <Stack key={i} direction={"row"} alignItems={"center"}>
                    <TextField
                      label={`opção ${i + 1}`}
                      required
                      fullWidth
                      value={o}
                      onChange={(e) => setOption(i, e.target.value)}
                      slotProps={{ htmlInput: { maxLength: 100 } }}
                    />
                    {options.length > MIN_OPTIONS && <Button onClick={() => removeOption(i)}>remover</Button>}
                  </Stack>
                ))}
              </Stack>
              {options.length < MAX_OPTIONS && (
                <Button onClick={() => setOptions((current) => [...current, ""])}>+ opção</Button>
              )}
            </Box>
            <TextField
              label={"encerrar em (opcional)"}
              type={"date"}
              value={closesOn}
              onChange={(e) => setClosesOn(e.target.value)}
              slotProps={{
                inputLabel: { shrink: true },
                htmlInput: { min: dayjs().format("YYYY-MM-DD") },
              }}
              helperText={"Sem data, a enquete fica aberta até você encerrá-la."}
            />
            {error && (
              <Typography variant={"body2"} color={"error"}>
                {error.message}
              </Typography>
            )}
            <Stack direction={"row"}>
              <Button type={"submit"} loading={isPending}>criar enquete</Button>
              <Button onClick={() => navigate(`/community/${communityId}/polls`)}>cancelar</Button>
            </Stack>
          </Stack>
        </form>
      </BoxWithTitle>
    </MembersOnly>
  );
};
