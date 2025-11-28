export interface Money {
    units: number;
    currency_code: string;
}

export interface Product {
    id: string;
    name: string;
    description: string;
    image_url: string;
    price: Money;
    stock?: number;
}

export interface ListProductResponse {
    total_pages: number;
    page: number;
    products: Product[];
}

export interface GetDetailsProducts {
    products: Product[];
}

