import React, { useState } from "react";
import SpeedDial from "@mui/material/SpeedDial";
import SpeedDialIcon from "@mui/material/SpeedDialIcon";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";

import Autocomplete from "@mui/material/Autocomplete";
import { useAuth } from "react-oidc-context";
import { EnvironmentData } from "../types/environment";
import { ApplicationData } from "../types/application";
import { LocationData } from "../types/location";
import { useSnackbar } from "notistack";
import { useGqlQuery, useQueryClient } from "../utils/http";
import gql from "graphql-tag";
import { useTranslation } from "react-i18next";

interface TFormData {
  environment: string;
  location: string;
  application: string;
  version: string;
}

const AddVersionDialog = ({ handleClose }: { handleClose: () => void }) => {
  const { t } = useTranslation();
  const defaultFormData = {
    environment: "",
    location: "",
    application: "",
    version: "",
  };
  const [formData, setFormData] = useState(defaultFormData);
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();

  const { data } = useGqlQuery(
    ["actionsform"],
    gql`
      query {
        environments {
          name
        }
        applications {
          name
        }
        locations {
          name
        }
      }
    `,
    {},
    (d) => ({
      environments: d.environments.map((r: EnvironmentData) => r.name),
      locations: d.locations.map((r: LocationData) => r.name),
      applications: d.applications.map((r: ApplicationData) => r.name),
    })
  );

  const auth = useAuth();

  const saveHandler = (formData: TFormData): Promise<Response> => {
    return fetch("/query", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${auth.user?.access_token}`,
      },
      body: JSON.stringify({
        query: `mutation { createVersion( input: { application:"${formData.application}", environment:"${formData.environment}", version:"${formData.version}", location:"${formData.location}" } ) { id } }`,
      }),
    });
  };

  const handleSubmit = (
    event: React.FormEvent<HTMLButtonElement | HTMLFormElement>
  ) => {
    event.preventDefault();
    saveHandler(formData)
      .then(() => {
        queryClient.refetchQueries({ exact: false, stale: true });
        enqueueSnackbar("New version created!", { variant: "success" });
        handleClose();
        setFormData(defaultFormData);
      })
      .catch((err) => {
        console.error(err);
        enqueueSnackbar("Failed to save!", { variant: "error" });
      });
  };

  return (
    <>
      <DialogTitle>{t("form_version_title")}</DialogTitle>
      <DialogContent>
        <DialogContentText>{t("form_version_description")}</DialogContentText>
        <form onSubmit={handleSubmit}>
          <Autocomplete
            disablePortal
            fullWidth
            freeSolo
            onInputChange={(event, value) => {
              setFormData({ ...formData, location: value });
            }}
            options={data?.locations || []}
            renderInput={(params) => (
              <TextField
                name="location"
                {...params}
                label={t("type_location")}
              />
            )}
          />
          <Autocomplete
            disablePortal
            fullWidth
            freeSolo
            onInputChange={(event, value) => {
              setFormData({ ...formData, environment: value });
            }}
            options={data?.environments || []}
            renderInput={(params) => (
              <TextField
                name="environment"
                {...params}
                label={t("type_environment")}
              />
            )}
          />
          <Autocomplete
            disablePortal
            fullWidth
            freeSolo
            onInputChange={(event, value) => {
              setFormData({ ...formData, application: value });
            }}
            options={data?.applications || []}
            renderInput={(params) => (
              <TextField
                name="application"
                {...params}
                label={t("type_application")}
              />
            )}
          />
          <TextField
            name="version"
            onChange={(event) => {
              setFormData({ ...formData, version: event.target.value });
            }}
            label={t("type_version")}
            fullWidth
          />
        </form>
      </DialogContent>
      <DialogActions>
        <Button variant="contained" onClick={handleClose}>
          {t("form_version_cancel")}
        </Button>
        <Button variant="contained" onClick={handleSubmit}>
          {t("form_version_submit")}
        </Button>
      </DialogActions>
    </>
  );
};

const Actions = () => {
  const [open, setOpen] = useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <SpeedDial
        ariaLabel="SpeedDial"
        sx={{ position: "fixed", bottom: 16, right: 16 }}
        icon={<SpeedDialIcon />}
        onClick={handleClickOpen}
      />
      <Dialog
        open={open}
        onClose={handleClose}
        sx={{
          "& .MuiTextField-root": { m: 1, ml: 0 },
          "& .MuiPaper-elevation": { overflow: "visible" },
        }}
      >
        {open && <AddVersionDialog handleClose={handleClose} />}
      </Dialog>
    </>
  );
};

export default Actions;
