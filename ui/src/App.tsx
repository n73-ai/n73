import { Routes, Route, BrowserRouter } from "react-router-dom";
import Landing from "./pages/landing";
import Layout from "./components/layout";
import { PrivateRoute } from "./lib/private-routes";
import NotFound from "./pages/not-found";
import Project from "./pages/project";
import Auth from "./pages/auth";

function App() {
  return (
    <BrowserRouter>
      <Routes>
      <Route path="signup" element={<Auth />} />
      <Route path="login" element={<Auth />} />
        <Route path="/" element={<Layout />}>
          <Route index element={<Landing />} />
        </Route>
        <Route path="/project/:projectID" element={<Project />} />
        <Route element={<PrivateRoute />}>
        </Route>
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
