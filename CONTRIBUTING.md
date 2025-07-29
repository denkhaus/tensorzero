# Contributing to TensorZero Go Client

We welcome contributions to the TensorZero Go Client! To ensure a smooth and efficient collaboration, please follow these guidelines.

## How to Contribute

1.  **Fork the repository:** Start by forking the `tensorzero-go` repository to your GitHub account.
2.  **Clone your fork:** Clone your forked repository to your local machine:

    ```bash
    git clone https://github.com/YOUR_USERNAME/tensorzero-go.git
    cd tensorzero-go
    ```

3.  **Create a new branch:** Create a new branch for your feature or bug fix:

    ```bash
    git checkout -b feature/your-feature-name
    # or
    git checkout -b bugfix/your-bug-fix-name
    ```

4.  **Make your changes:** Implement your feature or bug fix. Ensure your code adheres to the existing coding style and conventions.

5.  **Write tests:** If you're adding a new feature or fixing a bug, please write appropriate unit and/or integration tests to cover your changes. Refer to [TESTS.md](./docs/TESTS.md) for more details on testing.

6.  **Run tests:** Before committing, make sure all tests pass:

    ```bash
    make test-all
    ```

7.  **Commit your changes:** Write clear and concise commit messages. Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification if possible.

    ```bash
    git commit -m "feat: Add new awesome feature"
    # or
    git commit -m "fix: Fix a critical bug"
    ```

8.  **Push to your fork:** Push your changes to your forked repository:

    ```bash
    git push origin feature/your-feature-name
    ```

9.  **Create a Pull Request (PR):** Go to the original `tensorzero-go` repository on GitHub and create a new pull request from your forked branch. Please fill out the pull request template thoroughly.

## Code Style

We generally follow the standard Go formatting guidelines. Please run `go fmt` and `go vet` before committing your changes.

## Reporting Bugs

If you find a bug, please open an issue on GitHub using the "Bug report" template. Provide as much detail as possible, including steps to reproduce the bug, expected behavior, and your environment.

## Suggesting Enhancements

If you have an idea for a new feature or an improvement, please open an issue on GitHub using the "Feature request" template. Describe your idea clearly and explain why it would be beneficial.

## License

By contributing to this project, you agree that your contributions will be licensed under the Apache-2.0 License.
