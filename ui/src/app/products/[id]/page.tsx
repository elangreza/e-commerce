import ProductDescription from "@/components/ProductDescription";
import ProductTotalAndAddToCartProps from "@/components/ProductTotalAndAddToCart";
import { GetDetailsProducts } from "@/types/product";
import clsx from "clsx";
import Image from "next/image";

async function getProductsDetails(id: string): Promise<GetDetailsProducts> {
  const params = new URLSearchParams({
    "with_stock": "true",
    id,
  });

  const res = await fetch(`http://localhost:8080/api/product?${params.toString()}`)

  if (!res.ok) {
    throw new Error(`Failed to fetch products: ${res.status} ${res.statusText}`)
  }

  return res.json() as Promise<GetDetailsProducts>
}

interface PageProps {
  params: Promise<{
    id: string;
  }>
}

export default async function ProductsDetailPage({ params }: PageProps) {
  const { id } = await params;

  const products = await getProductsDetails(id)
  const product = products.products[0]

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-gray-600">
      <main className="min-h-screen w-full max-w-5xl flex-col items-center justify-between p-10 bg-white dark:bg-gray-500 sm:items-start">
        <div className="grid grid-rows-1 grid-cols-2 gap-12">
          <div
          >
            <Image
              src={product.image_url}
              width={400}
              height={400}
              alt={product.name}
              className="bg-white w-125 h-125 p-4 flex items-center justify-center"
            />
          </div>
          <div>
            <h1 className="text-2xl font-bold mb-4">{product.name}</h1>
            <p className={
              clsx(`my-2 font-bold`,
                product?.stock && product.stock > 3 ?
                  "text-green-400" :
                  product.stock === 0 ?
                    "text-red-400" : "text-yellow-400",
              )} >
              {product?.stock && product.stock > 3 ?
                "in stock" :
                product.stock === 0 ?
                  "out of stock" : `only ${product?.stock} left`}
            </p>
            <ProductDescription
              text={product.description}
            />
            <ProductTotalAndAddToCartProps
              stock={product?.stock || 0}
              price={product?.price}
              productID={product.id}
            />
          </div>
        </div>
      </main>
    </div>
  );
}
