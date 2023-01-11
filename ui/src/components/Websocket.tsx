import { useEffect } from "react";
import useWebSocket, { useEventSource } from "react-use-websocket";
import { useQueryClient } from "../utils/http";
import { useSnackbar } from "notistack";
import { useAuth } from "react-oidc-context";

const Websockets = () => {
  const authToken = useAuth().user?.access_token;

  let wsUrl = window.location.origin.replace("http", "ws") + "/socket";
  if (!process.env.NODE_ENV || process.env.NODE_ENV === "development") {
    wsUrl = "ws://127.0.0.1:8080/socket";
  }

  const queryClient = useQueryClient();
  const { enqueueSnackbar } = useSnackbar();

  const { lastMessage, sendMessage } = useWebSocket(wsUrl, {
    shouldReconnect: (closeEvent) => true,
    onOpen: (event) => {
      sendMessage(`Bearer ${authToken}`);
    },
  });

  let current: any;
  useEffect(() => {
    if (lastMessage && lastMessage !== current) {
      queryClient.refetchQueries({ exact: false, stale: true });
      enqueueSnackbar("New version received!", { variant: "info" });
      if (
        localStorage.getItem("setting.notification.enabled") === "granted" &&
        document.visibilityState !== "visible"
      ) {
        const notification = new Notification("Ledger", {
          body: "New version received!",
          icon: "/icon.png",
        });
        notification.onclick = () => {
          notification.close();
          window.parent.focus();
        };
      }
      current = lastMessage;
    }
  }, [lastMessage]);

  return <></>;
};

export default Websockets;
