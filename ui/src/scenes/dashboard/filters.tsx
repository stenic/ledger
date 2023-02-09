import {
  Autocomplete,
  Box,
  FormControl,
  FormGroup,
  IconButton,
  Popper,
  TextField,
} from "@mui/material";
import { useTranslation } from "react-i18next";
import { MouseEventHandler, useEffect, useState } from "react";
import FilterAltIcon from "@mui/icons-material/FilterAlt";

export const Filter = ({
  locations,
  environments,
  applications,
  initialState,
  filterCallback,
}: {
  locations: [string];
  environments: [string];
  applications: [string];
  initialState: {
    location: string;
    environment: string;
    application: string;
    day: string;
  };
  filterCallback?: (object: {
    location: string;
    environment: string;
    application: string;
    day: string;
  }) => void;
}) => {
  const { t } = useTranslation();
  const [filters, setFilters] = useState(initialState);

  const [hasFilters, setHasFilters] = useState(false);
  useEffect(() => {
    setHasFilters(
      filters.application !== "" ||
        filters.environment !== "" ||
        filters.location !== "" ||
        filters.day !== ""
    );
  }, [filters]);

  const [open, setOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);

  useEffect(() => {
    if (filterCallback) filterCallback(filters);
  }, [filters, filterCallback]);

  const canBeOpen = open && Boolean(anchorEl);
  const id = canBeOpen ? "filter-popper" : undefined;
  const handleFilterClick: MouseEventHandler<HTMLButtonElement> = (event) => {
    setAnchorEl(event.currentTarget);
    setOpen((previousOpen) => !previousOpen);
  };

  return (
    <>
      <IconButton aria-describedby={id} onClick={handleFilterClick}>
        <FilterAltIcon color={hasFilters ? "primary" : undefined} />
      </IconButton>
      <Popper
        id={id}
        open={open}
        anchorEl={anchorEl}
        onResize={undefined}
        onResizeCapture={undefined}
        placement={"bottom-end"}
      >
        <Box
          sx={{
            borderColor: "#000",
            borderRadius: "4px",
            bgcolor: "#434C5F",
            width: 300,
          }}
        >
          <Box
            sx={{
              display: "grid",
              gap: "10px",
              rowGap: "10px",
              gridTemplateColumns: "repeat(12, 1fr)",
              "& > div": {
                p: 2,
              },
            }}
          >
            <FormGroup
              sx={{
                gridColumn: "span 12",
              }}
            >
              <FormControl
                fullWidth
                sx={{
                  display: "grid",
                  gap: "10px",
                  rowGap: "10px",
                }}
              >
                <Autocomplete
                  disablePortal
                  fullWidth
                  freeSolo
                  onInputChange={(event, value) => {
                    setFilters({ ...filters, location: value });
                  }}
                  options={locations?.sort() || []}
                  value={filters.location}
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
                    setFilters({ ...filters, environment: value });
                  }}
                  options={environments?.sort() || []}
                  value={filters.environment}
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
                  value={filters.application}
                  onInputChange={(event, value) => {
                    setFilters({ ...filters, application: value });
                  }}
                  options={applications?.sort() || []}
                  renderInput={(params) => (
                    <TextField
                      name="application"
                      {...params}
                      label={t("type_application")}
                    />
                  )}
                />
                <TextField
                  label={t("type_date")}
                  type="date"
                  value={initialState.day}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  onChange={(event) => {
                    setFilters({ ...filters, day: event.target.value });
                  }}
                />
              </FormControl>
            </FormGroup>
          </Box>
        </Box>
      </Popper>
    </>
  );
};
