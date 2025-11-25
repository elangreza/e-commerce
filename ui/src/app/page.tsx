import Header from "@/components/Header";
import ProductCard from "@/components/ProductCard";
import { ListProductResponse } from "@/types/product";

async function getListProducts(page: number, search: string): Promise<ListProductResponse> {
  const params = new URLSearchParams({
    page: page.toString(), // Add page to URL parameters
    limit: '20',
    ...(search && { search })
  });

  const res = await fetch(`http://localhost:8080/api/products?${params.toString()}`)
  return res.json()
}

interface PageProps {
  searchParams: {
    search?: string;
    page?: string; // searchParams values are always strings
  }
}

export default async function Home({ searchParams }: PageProps) {

  // await is used to be getting the param in async way
  // see this link https://nextjs.org/docs/messages/sync-dynamic-apis#possible-ways-to-fix-it
  const { page, search } = await searchParams;

  const p = Number(page ?? '1'); // Convert string to number
  const s = search || "";
  const products = await getListProducts(p, s);

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-gray-600">
      <main className="min-h-screen w-full max-w-5xl flex-col items-center justify-between p-10 bg-white dark:bg-gray-500 sm:items-start">
        <Header />
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6">
          {products.products.map((product) => (
            <ProductCard product={product} key={product.id} />
          ))}
        </div>
      </main>
    </div>
  );
}
