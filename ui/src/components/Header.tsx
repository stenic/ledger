import { Typography, Box } from "@mui/material";
import { useEffect } from "react";
import { capitalizeFirstLetter } from "../utils/strings";

const Header = ({
  title,
  subtitle,
  hidden,
}: {
  title: string;
  subtitle?: string;
  hidden?: boolean;
}) => {
  useEffect(() => {
    document.title = title;
  }, [title, subtitle]);
  return (
    <Box mb="30px">
      <Typography
        variant="h2"
        fontWeight="bold"
        color="#e0e0e0"
        sx={{ m: "0 0 5px 0" }}
      >
        {title.toUpperCase()}
      </Typography>
      {subtitle && (
        <Typography variant="h5" color="#70d8bd">
          {capitalizeFirstLetter(subtitle)}
        </Typography>
      )}
    </Box>
  );
};

export default Header;
