import { Product } from "@/types/product"
import Image from "next/image"
import Link from "next/link"

interface ProductCardProps {
    product: Product
}

function ProductCard({ product }: ProductCardProps) {
    return (
        <Link href={`/products/${product.id}`} >
            <div className="border bg-white rounded-lg p-4 hover:shadow-md cursor-pointer">
                <Image
                    width={200}
                    height={200}
                    src={product.image_url}
                    alt={product.name}
                    className="w-full h-40 object-cover"

                />
                <h3 className="font-bold mt-2 text-black">{product.name}</h3>
                <p className="text-gray-600">{product.price.currency_code} - {product.price.units}</p>
            </div>
        </Link>
    )
}

export default ProductCard