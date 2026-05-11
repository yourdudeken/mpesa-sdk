from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="yourdudeken-mpesa-sdk",
    version="1.0.0",
    author="Kennedy Muthengi",
    author_email="ken@yourdudeken.com",
    description="A Python SDK for Mpesa Daraja APIs",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/yourdudeken/mpesa",
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
    ],
    packages=find_packages(),
    python_requires=">=3.8",
    install_requires=[
        "cryptography>=41.0.0",
        "requests>=2.28.0",
    ],
    keywords="mpesa daraja safaricom payments",
    project_urls={
        "Bug Tracker": "https://github.com/yourdudeken/mpesa/issues",
    },
)