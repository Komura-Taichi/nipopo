import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from "react-router";
import './index.css';
import Home from './pages/Home.tsx';
import RecordsList from "./pages/RecordsList.tsx";
import Contact from "./pages/Contact.tsx";
import RecordDetail from './pages/RecordDetail.tsx';
import Layout from "./layouts/Layout.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout />,
    children: [
      { index: true, Component: Home },
      { path: "records", Component: RecordsList },
      { path: "contact", Component: Contact },
      { path: "records/:recordId", Component: RecordDetail },
    ]
  },
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
);
