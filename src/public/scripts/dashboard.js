async function addProduct() {
    const url = document.getElementById("product-url").value;
    if (!url) return;
    await fetch("/api/product", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ "url":url }),
    });
    loadProducts();
}

async function loadProducts() {
    const response = await fetch("/api/products");
    const data = await response.json();
    const list = document.getElementById("product-list");
    list.innerHTML = "";
    data.tracked_products.forEach((p) => {
        const item = document.createElement("div");
        item.className = "product-item";
        item.innerHTML = `<img src="${p.ImgLink}" alt="${p.Name}"><span>${p.Name}</span>`;
        item.onclick = () => loadProductDetails(p.ID);
        list.appendChild(item);
    });
}

async function loadProductDetails(id) {
    const response = await fetch(`/api/product/${id}`, {
        method: "GET",
    });
    const data = await response.json();
    const details = document.getElementById("product-details");
    details.innerHTML = `<h2>${data.product_name}</h2><p>Current Price: ₹${data.current_price}</p>`;

    renderChart(data.price_history);
}

let priceChart = null; 

function renderChart(history) {
    const ctx = document.getElementById("priceChart").getContext("2d");

    
    if (priceChart) {
        priceChart.destroy();
    }

    
    priceChart = new Chart(ctx, {
        type: "line",
        data: {
            labels: history.map((h) => new Date(h.ChangedAt).toLocaleDateString()),
            datasets: [
                {
                    label: "Price",
                    data: history.map((h) => h.Price),
                    borderColor: "#007bff",
                    backgroundColor: "rgba(0, 123, 255, 0.1)",
                    fill: true,
                    tension: 0.4, 
                },
            ],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: true,
                    position: "top",
                },
            },
            scales: {
                x: {
                    title: { display: true, text: "Date" },
                },
                y: {
                    title: { display: true, text: "Price (₹)" },
                },
            },
        },
    });
}


loadProducts();