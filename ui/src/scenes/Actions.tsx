// import LedgerAdmin from "./admin";
import React, { useEffect, useState } from "react";
import SpeedDial from "@mui/material/SpeedDial";
import SpeedDialIcon from "@mui/material/SpeedDialIcon";
import SpeedDialAction from "@mui/material/SpeedDialAction";
import AddIcon from "@mui/icons-material/Add";
// import ShareIcon from "@mui/icons-material/Share";
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

interface TFormData {
  environment: string;
  application: string;
  version: string;
}

const AddVersionDialog = ({ handleClose }: { handleClose: () => void }) => {
  const defaultFormData = {
    environment: "",
    application: "",
    version: "",
  };
  const [formData, setFormData] = useState(defaultFormData);

  const auth = useAuth();
  const [apps, setApps] = useState([]);
  const [envs, setEnvs] = useState([]);

  useEffect(() => {
    const token = auth.user?.access_token;
    fetch("/query", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        query: "{ environments { name } }",
      }),
    })
      .then((e) => e.json())
      .then((d) => d.data.environments.map((r: EnvironmentData) => r.name))
      .then(setEnvs)
      .catch((e) => {
        console.error(e);
        // auth.signinSilent();
      });
  }, [auth, setEnvs]);

  useEffect(() => {
    const token = auth.user?.access_token;
    fetch("/query", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        query: "{ applications { name } }",
      }),
    })
      .then((e) => e.json())
      .then((d) => d.data.applications.map((r: ApplicationData) => r.name))
      .then(setApps)
      .catch((e) => {
        console.error(e);
        // auth.signinSilent();
      });
  }, [auth, setApps]);

  const saveHandler = (formData: TFormData) => {
    fetch("/query", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${auth.user?.access_token}`,
      },
      body: JSON.stringify({
        query: `mutation { createVersion( input: { application:"${formData.application}", environment:"${formData.environment}", version:"${formData.version}" } ) { id } }`,
      }),
    });
  };

  const handleSubmit = (
    event: React.FormEvent<HTMLButtonElement | HTMLFormElement>
  ) => {
    event.preventDefault();
    saveHandler(formData);
    handleClose();
    setFormData(defaultFormData);
  };

  return (
    <>
      <DialogTitle>Add version</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Fill out the form below to add a new version.
        </DialogContentText>
        <form onSubmit={handleSubmit}>
          <Autocomplete
            disablePortal
            fullWidth
            freeSolo
            onInputChange={(event, value) => {
              setFormData({ ...formData, environment: value });
            }}
            options={envs}
            renderInput={(params) => (
              <TextField name="environment" {...params} label="Environment" />
            )}
          />
          <Autocomplete
            disablePortal
            fullWidth
            freeSolo
            onInputChange={(event, value) => {
              setFormData({ ...formData, application: value });
            }}
            options={apps}
            renderInput={(params) => (
              <TextField name="application" {...params} label="Application" />
            )}
          />
          <TextField
            name="version"
            onChange={(event) => {
              setFormData({ ...formData, version: event.target.value });
            }}
            label="Version"
            fullWidth
          />
        </form>
      </DialogContent>
      <DialogActions>
        <Button variant="contained" onClick={handleClose}>
          Cancel
        </Button>
        <Button variant="contained" onClick={handleSubmit}>
          Add version
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

  const actions = [
    { icon: <AddIcon />, name: "Add new version", onClick: handleClickOpen },
    // { icon: <ShareIcon />, name: "Share" },
  ];

  return (
    <>
      <SpeedDial
        ariaLabel="SpeedDial"
        sx={{ position: "fixed", bottom: 16, right: 16 }}
        icon={<SpeedDialIcon />}
      >
        {actions.map((action) => (
          <SpeedDialAction
            key={action.name}
            icon={action.icon}
            tooltipTitle={action.name}
            onClick={action.onClick}
          />
        ))}
      </SpeedDial>
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
