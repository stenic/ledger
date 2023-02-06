import { Box, Typography } from "@mui/material";
import FlexBetween from "../../components/FlexBetween";

export const StatBox = ({
  title,
  value,
  icon,
  description,
}: {
  title: string;
  value: string;
  icon?: any;
  description: string;
}) => {
  return (
    <Box
      gridColumn="span 2"
      gridRow="span 1"
      display="flex"
      flexDirection="column"
      justifyContent="space-between"
      p="1.25 rem 1rem"
      flex="1 1 100%"
    >
      <FlexBetween>
        <Typography variant="h6">{title}</Typography>
        {icon}
      </FlexBetween>
      <Typography variant="h3" fontWeight="600">
        {value}
      </Typography>
      <Typography>{description}</Typography>
    </Box>
  );
};
