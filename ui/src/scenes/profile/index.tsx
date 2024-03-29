import { Box } from "@mui/material";
import Header from "../../components/Header";
import { useAuth } from "react-oidc-context";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import Typography from "@mui/material/Typography";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { useTranslation } from "react-i18next";

const Profile = () => {
  const auth = useAuth();
  const { t } = useTranslation();

  return (
    <Box m="20px">
      <Header
        title="Profile"
        subtitle={"Hello " + auth.user?.profile.preferred_username}
      />

      <Accordion defaultExpanded>
        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
          <Typography color="#4cceac" variant="h5">
            {t("profile_token")}
          </Typography>
        </AccordionSummary>
        <AccordionDetails>
          <code>
            <pre>{JSON.stringify(auth.user?.profile, null, 2)}</pre>
          </code>
        </AccordionDetails>
      </Accordion>
    </Box>
  );
};

export default Profile;
