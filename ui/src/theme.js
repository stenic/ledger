import { createTheme } from "@mui/material/styles";

const fontFamily = ["Source Sans Pro", "sans-serif"].join(",");

// color design tokens export
export const theme = createTheme({
  palette: {
    mode: "dark",
    primary: {
      main: "#4cceac",
    },
    secondary: {
      main: "#4cceac",
    },
    neutral: {
      dark: "#3d3d3d",
      main: "#666666",
      light: "#e0e0e0",
    },
    background: {
      default: "#141b2d",
      paper: "#1F2A40",
    },
    text: {
      primary: "#fff",
      secondary: "#ddd",
    },
  },
  typography: {
    fontFamily,
    fontSize: 12,
    h1: { fontFamily, fontSize: 40 },
    h2: { fontFamily, fontSize: 32 },
    h3: { fontFamily, fontSize: 24 },
    h4: { fontFamily, fontSize: 20 },
    h5: { fontFamily, fontSize: 16 },
    h6: { fontFamily, fontSize: 14 },
  },
});

export const nivoTheme = {
  tooltip: {
    container: {
      background: "#000",
      color: "#ffffff",
    },
  },
  axis: {
    ticks: {
      text: {
        fill: "#fff",
      },
    },
  },
  legends: {
    title: {
      text: {
        fill: "#fff",
      },
    },
    text: {
      fill: "#fff",
    },
    ticks: {
      text: {
        fill: "#fff",
      },
    },
  },
};

export const nivoBarProps = (data) => {
  return {
    groupMode: "stacked",
    margin: { top: 50, right: 130, bottom: 150, left: 60 },
    axisBottom: {
      tickRotation: 45,
    },
    tooltipLabel: (d) => d.id,
    enableLabel: false,
    legends: [
      {
        dataFrom: "keys",
        anchor: "bottom-right",
        direction: "column",
        justify: false,
        translateX: 120,
        translateY: 0,
        itemsSpacing: 2,
        itemWidth: 100,
        itemHeight: 20,
        itemDirection: "left-to-right",
        itemOpacity: 0.95,
        symbolSize: 20,
      },
    ],
  };
};
