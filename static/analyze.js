function clearPreviousResults() {
    const stockListDiv = document.getElementById('stockList');
    const topGainerLoserDiv = document.getElementById('topGainerLoser');
    const chartContainer = document.getElementById('chartContainer');
    stockListDiv.innerHTML = '';
    topGainerLoserDiv.innerHTML = '';
    chartContainer.innerHTML = ''; // Clear the chart as well
}

// Handle search
document.getElementById('searchForm').addEventListener('submit', (event) => {
    event.preventDefault();
    const symbol = document.getElementById('searchSymbol').value.toUpperCase();

    // Clear previous results before search
    clearPreviousResults();

    fetch(`/api/stocks/search?symbol=${symbol}`)
        .then(response => response.json())
        .then(data => {
            if (data) { // Check if data is returned
                const stockListDiv = document.getElementById('stockList');
                const resultDiv = document.createElement('div');
                resultDiv.innerHTML = ` 
                    <h2>Stock Today</h2>
                    <p>Symbol: ${data.symbol}</p>
                    <p>Price: ${data.price}</p>
                    <p>Change: ${data.change}</p>
                `;
                stockListDiv.appendChild(resultDiv);

            } else {
                alert('Stock not found!');
            }
        })
        .catch(() => alert('Stock not found!'));
});

// Sort stocks (Ascending)
document.getElementById('sortAsc').addEventListener('click', () => {
    // Clear previous results
    clearPreviousResults();

    fetch('/api/stocks/sort?order=asc')
        .then(response => response.json())
        .then(data => {
            const stockListDiv = document.getElementById('stockList');
            const table = document.createElement('table');
            const header = document.createElement('thead');
            header.innerHTML = ` 
                <tr>
                    <th>Symbol</th>
                    <th>Price</th>
                    <th>Change</th>
                </tr>
            `;
            table.appendChild(header);

            const tbody = document.createElement('tbody');
            data.forEach(stock => {
                const row = document.createElement('tr');
                row.innerHTML = ` 
                    <td>${stock.symbol}</td>
                    <td>${stock.price.toFixed(2)}</td>
                    <td>${stock.change}</td>
                `;
                tbody.appendChild(row);
            });
            table.appendChild(tbody);
            stockListDiv.appendChild(table); // Append sorted stock list

            // Create the bar chart with sorted data
            createBarChart(data);
        })
        .catch(() => console.error('Error sorting stocks.'));
});

// Sort stocks (Descending)
document.getElementById('sortDesc').addEventListener('click', () => {
    // Clear previous results
    clearPreviousResults();

    fetch('/api/stocks/sort?order=desc')
        .then(response => response.json())
        .then(data => {
            const stockListDiv = document.getElementById('stockList');
            const table = document.createElement('table');
            const header = document.createElement('thead');
            header.innerHTML = ` 
                <tr>
                    <th>Symbol</th>
                    <th>Price</th>
                    <th>Change</th>
                </tr>
            `;
            table.appendChild(header);

            const tbody = document.createElement('tbody');
            data.forEach(stock => {
                const row = document.createElement('tr');
                row.innerHTML = ` 
                    <td>${stock.symbol}</td>
                    <td>${stock.price.toFixed(2)}</td>
                    <td>${stock.change}</td>
                `;
                tbody.appendChild(row);
            });
            table.appendChild(tbody);
            stockListDiv.appendChild(table); // Append sorted stock list

            // Create the bar chart with sorted data
            createBarChart(data);
        })
        .catch(() => console.error('Error sorting stocks.'));
});

// Get top gainer and loser on button click
document.getElementById('topGainerLoserBtn').addEventListener('click', () => {
    // Clear previous results
    clearPreviousResults();

    fetch('/api/stocks/top-gainer-loser')
        .then(response => response.json())
        .then(data => {
            const resultDiv = document.getElementById('topGainerLoser');
            resultDiv.innerHTML = `
                <h2>Top Gainer</h2>
                <p>Symbol: ${data.top_gainer.symbol}, Price: ${data.top_gainer.price}, Change: ${data.top_gainer.change}</p>
                <h2>Top Loser</h2>
                <p>Symbol: ${data.top_loser.symbol}, Price: ${data.top_loser.price}, Change: ${data.top_loser.change}</p>
            `;
            // Create a bar chart for the top gainer and loser
            createBarChart([data.top_gainer, data.top_loser]);
        })
        .catch(() => console.error('Error fetching top gainer and loser.'));
});

function createBarChart(data) {
const width = 800;
const height = 400;
const margin = { top: 20, right: 30, bottom: 40, left: 40 };

// Define color scale: low price = red, high price = green
const colorScale = d3.scaleSequential(d3.interpolateRdYlGn)
.domain([0, d3.max(data, d => d.price)]); // Range from 0 to max price

const svg = d3.select('#chartContainer')
.append('svg')
.attr('width', width + margin.left + margin.right)
.attr('height', height + margin.top + margin.bottom)
.append('g')
.attr('transform', `translate(${margin.left}, ${margin.top})`);

const x = d3.scaleBand()
.domain(data.map(d => d.symbol))  // Use stock symbols for x-axis
.range([0, width])
.padding(0.1);

const y = d3.scaleLinear()
.domain([0, d3.max(data, d => d.price)])  // Use stock prices for y-axis
.nice()
.range([height, 0]);

// Clear previous bars before drawing new ones (important for dynamic updates)
svg.selectAll('.bar').remove();

// Add the tooltip div
const tooltip = d3.select('body').append('div')
.attr('class', 'tooltip')
.style('position', 'absolute')
.style('background-color', 'rgba(0, 0, 0, 0.7)')
.style('color', 'white')
.style('padding', '10px')
.style('border-radius', '4px')
.style('visibility', 'hidden');

// Draw bars with dynamic colors based on stock price
const bars = svg.append('g')
.selectAll('.bar')
.data(data)
.enter().append('rect')
.attr('class', 'bar')
.attr('x', d => x(d.symbol))
.attr('y', d => y(d.price))
.attr('width', x.bandwidth())
.attr('height', d => height - y(d.price))
.attr('fill', d => colorScale(d.price));
// Add event listeners for mouseover, mousemove, and mouseout
bars.on('mouseover', function (event, d) {
tooltip.style('visibility', 'visible')
    .html(`
        <strong>Symbol:</strong> ${d.symbol}<br>
        <strong>Price:</strong> $${d.price.toFixed(2)}<br>
        <strong>Change:</strong> ${d.change}%`)
    .style('left', (event.pageX + 10) + 'px')  // Position the tooltip to the right of the mouse
    .style('top', (event.pageY - 40) + 'px');  // Position the tooltip above the mouse

d3.select(this).attr('fill', '#ff7f0e');  // Highlight the bar on hover
})
.on('mousemove', function (event) {
tooltip.style('left', (event.pageX + 10) + 'px')  // Update the tooltip position
    .style('top', (event.pageY - 40) + 'px');
})
.on('mouseout', function () {
tooltip.style('visibility', 'hidden');  // Hide the tooltip
d3.select(this).attr('fill', d => colorScale(d.price));  // Reset the bar color
});

// X Axis
svg.append('g')
.attr('class', 'x-axis')
.attr('transform', `translate(0,${height})`)
.call(d3.axisBottom(x));

// Y Axis
svg.append('g')
.attr('class', 'y-axis')
.call(d3.axisLeft(y));
}