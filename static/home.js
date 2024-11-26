// Populate the stock table
function loadStocks() {
    fetch('/api/stocks')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('stocksTableBody');
            tableBody.innerHTML = '';
            data.forEach(stock => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${stock.symbol}</td>
                    <td>${stock.price.toFixed(2)}</td>
                    <td>${stock.change}</td>
                `;
                tableBody.appendChild(row);
            });
        })
        .catch(error => console.error('Error fetching stocks:', error));
}

// Start scraping
document.getElementById('scrapeBtn').addEventListener('click', () => {
    fetch('/api/stocks/scrape')
        .then(response => response.json())
        .then(() => loadStocks()) // Reload table after scraping
        .catch(error => console.error('Error scraping stocks:', error));
});

// Load stocks on page load
window.onload = loadStocks;
