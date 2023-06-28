import { Layout } from "components/Layout";
import { Link } from "react-router-dom";

export function NotFound() {
  return (
    <Layout>
      <section className="w-[100vw] h-[100vh] flex flex-col justify-center items-center gap-3">
        <h1 className="text-6xl font-semibold">Whoops!</h1>

        <p>
          Some pages are so highly classified that even we are unaware of their
          existence.
        </p>

        <Link to="/" className="mt-5">
          Go back Home
        </Link>
      </section>
    </Layout>
  );
}
