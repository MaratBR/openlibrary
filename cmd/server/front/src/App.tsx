import { RouterProvider } from "react-router-dom";
import router from "./router";
import { QueryClientProvider } from "@tanstack/react-query";
import queryClient from "./query-client";
import { TooltipProvider } from "./components/ui/tooltip";

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <RouterProvider router={router} />
      </TooltipProvider>
    </QueryClientProvider>
  );
}

export default App;
