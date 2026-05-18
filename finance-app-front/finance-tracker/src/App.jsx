import React, { useState, useEffect } from 'react';

const API_BASE = 'http://localhost:8080';

const API = {
  auth: {
    login: `${API_BASE}/auth/login`,
    register: `${API_BASE}/auth/register`,
  },
  me: `${API_BASE}/api/me`,
  transactions: `${API_BASE}/api/transactions`,
  categories: `${API_BASE}/api/categories`,
  budgets: `${API_BASE}/api/budgets`,
  summary: `${API_BASE}/api/analytics/summary`,
};

export default function FinanceApp() {
  const [authToken, setAuthToken] = useState(localStorage.getItem('token'));
  const [currentPage, setCurrentPage] = useState('dashboard');
  const [summary, setSummary] = useState(null);
  const [transactions, setTransactions] = useState([]);
  const [categories, setCategories] = useState([]);
  const [budgets, setBudgets] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (authToken) {
      loadData();
      const interval = setInterval(loadData, 30000);
      return () => clearInterval(interval);
    }
  }, [authToken]);

  const loadData = async () => {
    try {
      const headers = { 'Authorization': `Bearer ${authToken}` };

      const [summaryRes, txRes, catRes, budgetRes] = await Promise.all([
        fetch(API.summary, { headers }),
        fetch(API.transactions, { headers }),
        fetch(API.categories, { headers }),
        fetch(API.budgets, { headers })
      ]);

      if (summaryRes.ok) setSummary(await summaryRes.json());
      if (txRes.ok) setTransactions(await txRes.json());
      if (catRes.ok) setCategories(await catRes.json());
      if (budgetRes.ok) {
  const data = await budgetRes.json();

  setBudgets(
    data.map(b => ({
      id: b.ID,
      category_id: b.CategoryID,
      limit: Number(b.Limit),
      period: b.Period,
      createdAt: b.CreatedAt,
    }))
  );
}
    } catch (err) {
      console.error('Ошибка загрузки данных:', err);
    }
  };

  const handleLogout = () => {
    setAuthToken(null);
    localStorage.removeItem('token');
    setCurrentPage('login');
  };

  if (!authToken) {
    return <AuthPage setAuthToken={setAuthToken} setError={setError} error={error} />;
  }

  return (
    <div style={styles.container}>
      <Sidebar currentPage={currentPage} setCurrentPage={setCurrentPage} onLogout={handleLogout} />

      <main style={styles.main}>
        {currentPage === 'dashboard' && <Dashboard summary={summary} transactions={transactions} />}
        {currentPage === 'transactions' && <TransactionsPage transactions={transactions} categories={categories} onRefresh={loadData} />}
        {currentPage === 'categories' && <CategoriesPage categories={categories} onRefresh={loadData} />}
        {currentPage === 'budgets' && <BudgetsPage budgets={budgets} categories={categories} onRefresh={loadData} />}
      </main>
    </div>
  );
}

function AuthPage({ setAuthToken, setError, error }) {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const endpoint = isLogin ? API.auth.login : API.auth.register;

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      if (!response.ok) throw new Error('Ошибка авторизации');

      const data = await response.json();
      localStorage.setItem('token', data.token);
      setAuthToken(data.token);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.authContainer}>
      <div style={styles.authCard}>
        <div style={styles.authHeader}>
          <i className="ti ti-wallet" style={styles.authIcon}></i>
          <h1 style={styles.authTitle}>Финансовый трекер</h1>
        </div>

        <form onSubmit={handleSubmit} style={styles.form}>
          <input
            type="email"
            placeholder="Электронная почта"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            style={styles.input}
            required
          />
          <input
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={styles.input}
            required
          />

          {error && <div style={styles.error}>{error}</div>}

          <button type="submit" style={styles.primaryBtn} disabled={loading}>
            {loading ? 'Загрузка...' : isLogin ? 'Войти' : 'Создать аккаунт'}
          </button>
        </form>

        <button
          onClick={() => setIsLogin(!isLogin)}
          style={styles.toggleBtn}
        >
          {isLogin ? 'Нет аккаунта?' : 'Уже есть аккаунт?'}
        </button>
      </div>
    </div>
  );
}

function Sidebar({ currentPage, setCurrentPage, onLogout }) {
  const menuItems = [
    { id: 'dashboard', label: 'Панель управления', icon: 'ti-layout-dashboard' },
    { id: 'transactions', label: 'Транзакции', icon: 'ti-cash' },
    { id: 'categories', label: 'Категории', icon: 'ti-tags' },
    { id: 'budgets', label: 'Бюджеты', icon: 'ti-chart-bar' }
  ];

  return (
    <aside style={styles.sidebar}>
      <div style={styles.sidebarHeader}>
        <i className="ti ti-wallet" style={styles.sidebarIcon}></i>
        <h2 style={styles.sidebarTitle}>FinTrack</h2>
      </div>

      <nav style={styles.nav}>
        {menuItems.map(item => (
          <button
            key={item.id}
            onClick={() => setCurrentPage(item.id)}
            style={{
              ...styles.navItem,
              ...(currentPage === item.id ? styles.navItemActive : {})
            }}
          >
            <i className={`ti ${item.icon}`} style={styles.navIcon}></i>
            {item.label}
          </button>
        ))}
      </nav>

      <button onClick={onLogout} style={styles.logoutBtn}>
        <i className="ti ti-logout" style={styles.navIcon}></i>
        Выйти
      </button>
    </aside>
  );
}

function Dashboard({ summary, transactions }) {
  const recentTx = transactions?.slice(0, 5) || [];

  return (
    <div style={styles.pageContent}>
      <h1 style={styles.pageTitle}>Панель управления</h1>

      {summary && (
        <div style={styles.summaryGrid}>
          <SummaryCard title="Общий доход" amount={summary.total_income} type="income" icon="ti-arrow-down" />
          <SummaryCard title="Общие расходы" amount={summary.total_expense} type="expense" icon="ti-arrow-up" />
          <SummaryCard title="Баланс" amount={summary.balance} type="balance" icon="ti-wallet" />
        </div>
      )}

      <div style={styles.section}>
        <h2 style={styles.sectionTitle}>Последние транзакции</h2>
        {recentTx.length > 0 ? (
          <div style={styles.transactionList}>
            {recentTx.map(tx => (
              <TransactionItem key={tx.id} transaction={tx} />
            ))}
          </div>
        ) : (
          <p style={styles.emptyState}>Транзакций пока нет</p>
        )}
      </div>
    </div>
  );
}

function SummaryCard({ title, amount, type, icon }) {
  const typeColors = {
    income: { bg: 'var(--color-background-success)', text: 'var(--color-text-success)' },
    expense: { bg: 'var(--color-background-danger)', text: 'var(--color-text-danger)' },
    balance: { bg: 'var(--color-background-info)', text: 'var(--color-text-info)' }
  };

  const colors = typeColors[type];

  return (
    <div style={{ ...styles.card, backgroundColor: colors.bg }}>
      <div style={styles.cardContent}>
        <p style={styles.cardLabel}>{title}</p>
        <div style={styles.cardAmount}>
          <i className={`ti ${icon}`} style={{ ...styles.cardIcon, color: colors.text }}></i>
          <h3 style={{ ...styles.cardValue, color: colors.text }}>
            {typeof amount === 'string' ? amount : amount?.toFixed(2) || '0.00'}
          </h3>
        </div>
      </div>
    </div>
  );
}

function TransactionItem({ transaction }) {
  const isIncome = transaction.TransactionType === 'IN';
  const color = isIncome ? 'var(--color-text-success)' : 'var(--color-text-danger)';
  const sign = isIncome ? '+' : '-';

  return (
    <div style={styles.transactionRow}>
      <div style={styles.txInfo}>
        <p style={styles.txCategory}>{transaction.Category}</p>
        <p style={styles.txDescription}>{transaction.Description || 'Без описания'}</p>
      </div>

      <p style={{ ...styles.txAmount, color }}>
        {sign}{transaction.Amount}
      </p>
    </div>
  );
}

function TransactionsPage({ transactions, categories, onRefresh }) {
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    amount: '',
    type: 'OUT',
    category: '',
    description: ''
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(API.transactions, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          ...formData,
          amount: parseFloat(formData.amount)
        })
      });

      if (response.ok) {
        setShowForm(false);
        setFormData({ amount: '', type: 'OUT', category: '', description: '' });
        onRefresh();
      }
    } catch (err) {
      console.error('Ошибка создания транзакции:', err);
    }
  };

  return (
    <div style={styles.pageContent}>
      <div style={styles.pageHeader}>
        <h1 style={styles.pageTitle}>Транзакции</h1>
        <button onClick={() => setShowForm(!showForm)} style={styles.primaryBtn}>
          <i className="ti ti-plus" style={styles.btnIcon}></i> Добавить транзакцию
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={styles.formCard}>
          <div style={styles.formGrid}>
            <input type="number" step="0.01" placeholder="Сумма" value={formData.amount}
              onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
              style={styles.input} required />

            <select value={formData.type}
              onChange={(e) => setFormData({ ...formData, type: e.target.value })}
              style={styles.select}>
              <option value="OUT">Расход</option>
              <option value="IN">Доход</option>
            </select>

            <input type="text" placeholder="Категория" value={formData.category}
              onChange={(e) => setFormData({ ...formData, category: e.target.value })}
              style={styles.input} required />

            <input type="text" placeholder="Описание (необязательно)" value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              style={styles.input} />
          </div>

          <div style={styles.formActions}>
            <button type="submit" style={styles.primaryBtn}>Сохранить транзакцию</button>
            <button type="button" onClick={() => setShowForm(false)} style={styles.secondaryBtn}>Отмена</button>
          </div>
        </form>
      )}

      {transactions?.length > 0 ? (
        <div style={styles.transactionList}>
          {transactions.map(tx => (
            <TransactionItem key={tx.id} transaction={tx} />
          ))}
        </div>
      ) : (
        <p style={styles.emptyState}>Транзакций пока нет. Создайте первую!</p>
      )}
    </div>
  );
}

function CategoriesPage({ categories, onRefresh }) {
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({ name: '', type: 'OUT' });

  const handleSubmit = async (e) => {
    e.preventDefault();

    await fetch(API.categories, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        Name: formData.name,
        TransactionType: formData.type
      })
    });

    setShowForm(false);
    setFormData({ name: '', type: 'OUT' });
    onRefresh();
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Удалить категорию?')) return;

    await fetch(`${API_BASE}/api/categories/${id}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    onRefresh();
  };

  return (
    <div style={styles.pageContent}>
      <div style={styles.pageHeader}>
        <h1 style={styles.pageTitle}>Категории</h1>
        <button onClick={() => setShowForm(!showForm)} style={styles.primaryBtn}>
          Добавить
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={styles.formCard}>
          <input
            placeholder="Название"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            style={styles.input}
            required
          />

          <select
            value={formData.type}
            onChange={(e) => setFormData({ ...formData, type: e.target.value })}
            style={styles.select}
          >
            <option value="OUT">Расход</option>
            <option value="IN">Доход</option>
          </select>

          <button type="submit" style={styles.primaryBtn}>
            Сохранить
          </button>
        </form>
      )}

      <div style={styles.categoriesGrid}>
        {categories?.map(cat => (
          <div key={cat.ID} style={styles.categoryCard}>
            <div style={styles.categoryInfo}>
              <p style={styles.categoryName}>{cat.Name}</p>

              <span style={{
                ...styles.badge,
                backgroundColor:
                  cat.TransactionType === 'IN'
                    ? 'var(--color-background-success)'
                    : 'var(--color-background-danger)'
              }}>
                {cat.TransactionType === 'IN' ? 'Доход' : 'Расход'}
              </span>
            </div>

            <button
              onClick={() => handleDelete(cat.ID)}
              style={styles.deleteBtn}
            >
              🗑
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

function BudgetsPage({ budgets, categories, onRefresh }) {
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({ category_id: '', limit: '', period: 'MONTH' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(API.budgets, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          ...formData,
          limit: parseFloat(formData.limit)
        })
      });

      if (response.ok) {
        setShowForm(false);
        setFormData({ category_id: '', limit: '', period: 'MONTH' });
        onRefresh();
      }
    } catch (err) {
      console.error('Ошибка создания бюджета:', err);
    }
  };

const getCategoryName = (id) => {
  return categories?.find(cat => cat.ID === id)?.Name || 'Неизвестно';
};

  return (
    <div style={styles.pageContent}>
      <div style={styles.pageHeader}>
        <h1 style={styles.pageTitle}>Бюджеты</h1>
        <button onClick={() => setShowForm(!showForm)} style={styles.primaryBtn}>
          <i className="ti ti-plus" style={styles.btnIcon}></i> Установить бюджет
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={styles.formCard}>
          <div style={styles.formGrid}>
            <select
              value={formData.category_id}
              onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
              style={styles.select}
              required
            >
              <option value="">Выберите категорию</option>
   {categories?.map(cat => (
  <option key={cat.ID} value={cat.ID}>
    {cat.Name}
  </option>
))}
            </select>

            <input
              type="number"
              step="0.01"
              placeholder="Лимит суммы"
              value={formData.limit}
              onChange={(e) => setFormData({ ...formData, limit: e.target.value })}
              style={styles.input}
              required
            />

            <select
              value={formData.period}
              onChange={(e) => setFormData({ ...formData, period: e.target.value })}
              style={styles.select}
            >
              <option value="MONTH">Ежемесячно</option>
              <option value="WEEK">Еженедельно</option>
              <option value="YEAR">Ежегодно</option>
            </select>
          </div>

          <div style={styles.formActions}>
            <button type="submit" style={styles.primaryBtn}>Создать бюджет</button>
            <button type="button" onClick={() => setShowForm(false)} style={styles.secondaryBtn}>Отмена</button>
          </div>
        </form>
      )}

      <div style={styles.budgetsList}>
        {budgets?.length > 0 ? (
          budgets.map(budget => (
            <div key={budget.id} style={styles.budgetCard}>
              <div>
                <p style={styles.budgetCategory}>{getCategoryName(budget.category_id)}</p>
                <p style={styles.budgetPeriod}>{budget.period}</p>
              </div>
              <p style={styles.budgetLimit}>{budget.limit}</p>
            </div>
          ))
        ) : (
          <p style={styles.emptyState}>Бюджеты пока не заданы</p>
        )}
      </div>
    </div>
  );
}

const styles = {
  container: {
    display: 'flex',
    height: '100vh',
    fontFamily: 'var(--font-sans)',
    backgroundColor: 'var(--color-background-tertiary)',
    color: 'var(--color-text-primary)'
  },
  sidebar: {
    width: '240px',
    backgroundColor: 'var(--color-background-primary)',
    borderRight: '1px solid var(--color-border-tertiary)',
    padding: '2rem 1rem',
    display: 'flex',
    flexDirection: 'column',
    overflowY: 'auto'
  },
  sidebarHeader: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.75rem',
    marginBottom: '2rem'
  },
  sidebarIcon: {
    fontSize: '24px',
    color: 'var(--color-text-info)'
  },
  sidebarTitle: {
    fontSize: '18px',
    fontWeight: '500',
    margin: '0',
    color: 'var(--color-text-primary)'
  },
  nav: {
    display: 'flex',
    flexDirection: 'column',
    gap: '0.5rem',
    flex: '1'
  },
  navItem: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.75rem',
    padding: '0.75rem 1rem',
    backgroundColor: 'transparent',
    border: 'none',
    borderRadius: 'var(--border-radius-md)',
    cursor: 'pointer',
    color: 'var(--color-text-secondary)',
    fontSize: '14px',
    fontWeight: '400',
    transition: 'all 0.2s'
  },
  navItemActive: {
    backgroundColor: 'var(--color-background-secondary)',
    color: 'var(--color-text-info)',
    fontWeight: '500'
  },
  navIcon: {
    fontSize: '18px'
  },
  logoutBtn: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.75rem',
    padding: '0.75rem 1rem',
    backgroundColor: 'var(--color-background-danger)',
    border: 'none',
    borderRadius: 'var(--border-radius-md)',
    cursor: 'pointer',
    color: 'var(--color-text-danger)',
    fontSize: '14px',
    fontWeight: '500',
    marginTop: 'auto'
  },
  main: {
    flex: '1',
    overflowY: 'auto',
    padding: '2rem'
  },
  pageContent: {
    maxWidth: '1200px',
    margin: '0 auto'
  },
  pageHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: '2rem'
  },
  pageTitle: {
    fontSize: '22px',
    fontWeight: '500',
    margin: '0',
    color: 'var(--color-text-primary)'
  },
  summaryGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
    gap: '1.5rem',
    marginBottom: '2rem'
  },
card: {
  backgroundColor: 'var(--color-background-primary)',
  borderRadius: 'var(--border-radius-lg)',
  padding: '1.5rem',
  border: '1px solid var(--color-border-tertiary)',
  boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
},
  cardContent: {
    display: 'flex',
    flexDirection: 'column',
    gap: '1rem'
  },
  cardLabel: {
    fontSize: '13px',
    color: 'var(--color-text-secondary)',
    margin: '0',
    fontWeight: '400'
  },
  cardAmount: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.75rem'
  },
  cardIcon: {
    fontSize: '24px'
  },
  cardValue: {
    fontSize: '20px',
    fontWeight: '500',
    margin: '0'
  },
  section: {
    marginBottom: '2rem'
  },
  sectionTitle: {
    fontSize: '18px',
    fontWeight: '500',
    margin: '0 0 1rem 0',
    color: 'var(--color-text-primary)'
  },
  transactionList: {
    backgroundColor: 'var(--color-background-primary)',
    borderRadius: 'var(--border-radius-lg)',
    border: '1px solid var(--color-border-tertiary)',
    overflow: 'hidden'
  },
  transactionRow: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: '1rem 1.5rem',
    borderBottom: '1px solid var(--color-border-tertiary)',
    '&:last-child': { borderBottom: 'none' }
  },
  txInfo: {
    flex: '1'
  },
  txCategory: {
    fontSize: '14px',
    fontWeight: '500',
    margin: '0 0 0.25rem 0',
    color: 'var(--color-text-primary)'
  },
  txDescription: {
    fontSize: '13px',
    color: 'var(--color-text-secondary)',
    margin: '0'
  },
  txAmount: {
    fontSize: '16px',
    fontWeight: '500',
    margin: '0'
  },
  emptyState: {
    textAlign: 'center',
    color: 'var(--color-text-secondary)',
    padding: '2rem',
    fontSize: '14px'
  },
  formCard: {
    backgroundColor: 'var(--color-background-primary)',
    borderRadius: 'var(--border-radius-lg)',
    border: '1px solid var(--color-border-tertiary)',
    padding: '1.5rem',
    marginBottom: '2rem'
  },
  formGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
    gap: '1rem',
    marginBottom: '1.5rem'
  },
  input: {
    padding: '0.75rem',
    border: '1px solid var(--color-border-secondary)',
    borderRadius: 'var(--border-radius-md)',
    fontSize: '14px',
    fontFamily: 'var(--font-sans)',
    backgroundColor: 'var(--color-background-primary)',
    color: 'var(--color-text-primary)'
  },
  select: {
    padding: '0.75rem',
    border: '1px solid var(--color-border-secondary)',
    borderRadius: 'var(--border-radius-md)',
    fontSize: '14px',
    fontFamily: 'var(--font-sans)',
    backgroundColor: 'var(--color-background-primary)',
    color: 'var(--color-text-primary)'
  },
  formActions: {
    display: 'flex',
    gap: '1rem'
  },
  primaryBtn: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.5rem',
    padding: '0.75rem 1.5rem',
    backgroundColor: 'var(--color-background-info)',
    border: 'none',
    borderRadius: 'var(--border-radius-md)',
    color: 'var(--color-text-info)',
    fontWeight: '500',
    cursor: 'pointer',
    fontSize: '14px',
    transition: 'all 0.2s'
  },
  secondaryBtn: {
    padding: '0.75rem 1.5rem',
    backgroundColor: 'var(--color-background-secondary)',
    border: 'none',
    borderRadius: 'var(--border-radius-md)',
    color: 'var(--color-text-secondary)',
    fontWeight: '500',
    cursor: 'pointer',
    fontSize: '14px'
  },
  btnIcon: {
    fontSize: '16px'
  },
  categoriesGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))',
    gap: '1rem'
  },
  categoryCard: {
    backgroundColor: 'var(--color-background-primary)',
    borderRadius: 'var(--border-radius-lg)',
    border: '1px solid var(--color-border-tertiary)',
    padding: '1.5rem',
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center'
  },
  categoryInfo: {
    flex: '1'
  },
  categoryName: {
    fontSize: '16px',
    fontWeight: '500',
    margin: '0 0 0.5rem 0',
    color: 'var(--color-text-primary)'
  },
  badge: {
    display: 'inline-block',
    padding: '0.25rem 0.75rem',
    borderRadius: 'var(--border-radius-md)',
    fontSize: '12px',
    fontWeight: '500',
    opacity: 0.7
  },
  deleteBtn: {
    backgroundColor: 'transparent',
    border: 'none',
    color: 'var(--color-text-danger)',
    cursor: 'pointer',
    fontSize: '16px',
    padding: '0.5rem'
  },
  budgetsList: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))',
    gap: '1rem'
  },
  budgetCard: {
    backgroundColor: 'var(--color-background-primary)',
    borderRadius: 'var(--border-radius-lg)',
    border: '1px solid var(--color-border-tertiary)',
    padding: '1.5rem',
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center'
  },
  budgetCategory: {
    fontSize: '16px',
    fontWeight: '500',
    margin: '0 0 0.5rem 0',
    color: 'var(--color-text-primary)'
  },
  budgetPeriod: {
    fontSize: '13px',
    color: 'var(--color-text-secondary)',
    margin: '0'
  },
  budgetLimit: {
    fontSize: '18px',
    fontWeight: '500',
    color: 'var(--color-text-info)',
    margin: '0'
  },
  authContainer: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    minHeight: '100vh',
    backgroundColor: 'var(--color-background-tertiary)',
    padding: '1rem'
  },
  authCard: {
    backgroundColor: 'var(--color-background-primary)',
    borderRadius: 'var(--border-radius-lg)',
    border: '1px solid var(--color-border-tertiary)',
    padding: '2rem',
    width: '100%',
    maxWidth: '400px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.05)'
  },
  authHeader: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: '1rem',
    marginBottom: '2rem'
  },
  authIcon: {
    fontSize: '36px',
    color: 'var(--color-text-info)'
  },
  authTitle: {
    fontSize: '22px',
    fontWeight: '500',
    margin: '0',
    color: 'var(--color-text-primary)'
  },
  form: {
    display: 'flex',
    flexDirection: 'column',
    gap: '1rem',
    marginBottom: '1.5rem'
  },
  error: {
    backgroundColor: 'var(--color-background-danger)',
    color: 'var(--color-text-danger)',
    padding: '0.75rem',
    borderRadius: 'var(--border-radius-md)',
    fontSize: '13px',
    opacity: 0.7
  },
  toggleBtn: {
    backgroundColor: 'transparent',
    border: 'none',
    color: 'var(--color-text-info)',
    cursor: 'pointer',
    fontSize: '14px',
    fontWeight: '500',
    textAlign: 'center',
    width: '100%'
  }
};