import ProductCard from "@/components/ProductCard";
import { ListProductResponse } from "@/types/product";

async function getListProducts(page: number, search: string): Promise<ListProductResponse> {
  const params = new URLSearchParams({
    page: page.toString(), // Add page to URL parameters
    limit: '20',
    ...(search && { search })
  });

  const res = await fetch(`http://localhost:8080/api/products?${params.toString()}`)
  if (!res.ok) {
    throw new Error(`Failed to fetch products: ${res.status} ${res.statusText}`)
  }

  return res.json() as Promise<ListProductResponse>
}

interface PageProps {
  searchParams: Promise<{
    search?: string;
    page?: string; // searchParams values are always strings
  }>
}

export default async function Home({ searchParams }: PageProps) {

  // await is used to be getting the param in async way
  // see this link https://nextjs.org/docs/messages/sync-dynamic-apis#possible-ways-to-fix-it
  const { page, search } = await searchParams;

  const p = Number(page ?? '1'); // Convert string to number
  const s = search || "";
  let products: ListProductResponse | null = null;
  let errorMsg: string | null = null;

  try {
    products = await getListProducts(p, s);
  } catch (err) {
    errorMsg = err instanceof Error ? err.message : String(err);
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-gray-600">
      <main className="min-h-screen w-full max-w-5xl flex-col items-center justify-between p-10 bg-white dark:bg-gray-500 sm:items-start">
        {errorMsg ? (
          <div className="w-full p-6">
            <h2 className="text-xl font-semibold text-red-600">Error loading products</h2>
            <p className="mt-2 text-sm text-gray-700">{errorMsg}</p>
            <a
              href={`/?search=${encodeURIComponent(s)}&page=${p}`}
              className="mt-3 inline-block text-sm text-blue-600"
            >
              Retry
            </a>
          </div>
        ) : products?.products == null || products?.products?.length === 0 ? (
          <div className="w-full p-24 text-center">
            <h2 className="text-xl font-semibold">No products found</h2>
            <p className="mt-2">Try a different search or clear filters.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-6">
            {(products?.products ?? []).map((product) => (
              <ProductCard product={product} key={product.id} />
            ))}
          </div>
        )}
      </main>
    </div>

  );
}
